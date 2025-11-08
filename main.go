package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"email-api/mail"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: no se cargó el archivo .env, utilizando variables de entorno del sistema.")
	}
}

// Estructura para los productos recomendados
type Product struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	BuyURL      string `json:"buy_url"`
}

// Estructura para la solicitud de recomendaciones
type RecommendationRequest struct {
	UserName         string    `json:"user_name"`
	Subject          string    `json:"subject"`
	Products         []Product `json:"products"`
	CallToActionURL  string    `json:"call_to_action_url"`
	PhoneNumber      string    `json:"phone_number"`
	DestinationEmail string    `json:"destination_email"`
}

// Estructura para la llamada telefónica
type PhoneCallRequest struct {
	PhoneNumber string `json:"phone_number"`
}

// Estructura para recibir los datos del correo (mantenida para compatibilidad)
type EmailRequest struct {
	Mail    string `json:"mail"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// Habilitar CORS
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// Generar la sección de productos en HTML
func generateProductsSection(products []Product) string {
	var productsHTML string
	for _, product := range products {
		imageSection := fmt.Sprintf(`<div class="product-image">%s</div>`, product.Image)
		if product.Image != "" {
			imageSection = fmt.Sprintf(`<img src="%s" alt="%s" style="width: 100%%; height: 150px; object-fit: cover; border-radius: 6px; margin-bottom: 15px;">`, product.Image, product.Name)
		}

		productsHTML += fmt.Sprintf(`
			<div class="recommendation-section">
				<div class="recommendation-card">
					%s
					<h3>%s</h3>
					<p>%s</p>
					<a href="%s" class="buy-btn">BUY NOW</a>
				</div>
			</div>
		`, imageSection, product.Name, product.Description, product.BuyURL)
	}
	return productsHTML
}

// Generar el HTML completo de recomendaciones
func generateRecommendationHTML(req RecommendationRequest) string {
	productsSection := generateProductsSection(req.Products)

	htmlTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Product Recommendations</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            line-height: 1.6;
            color: #000;
            margin: 0;
            padding: 0;
            background-color: #f0f0f0;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
        }
        .header-banner {
            background: #000000;
            padding: 40px 20px;
            text-align: center;
            color: white;
        }
        .header-banner h1 {
            margin: 0;
            font-size: 28px;
            font-weight: 700;
        }
        .content {
            padding: 40px 20px;
        }
        .greeting {
            font-size: 24px;
            font-weight: 600;
            margin-bottom: 20px;
            color: #000;
        }
        .intro-text {
            font-size: 14px;
            line-height: 1.8;
            color: #333;
            margin-bottom: 30px;
        }
        .recommendation-section {
            margin-bottom: 40px;
        }
        .recommendation-card {
            border: 2px solid #000;
            border-radius: 8px;
            padding: 25px;
            text-align: center;
            background-color: #fafafa;
            margin-bottom: 15px;
        }
        .recommendation-card h3 {
            margin: 0 0 15px 0;
            font-size: 18px;
            font-weight: 600;
            color: #000;
        }
        .recommendation-card p {
            margin: 0 0 15px 0;
            font-size: 13px;
            color: #333;
        }
        .product-image {
            width: 100%%;
            height: 150px;
            background: #cccccc;
            border-radius: 6px;
            margin-bottom: 15px;
            display: flex;
            align-items: center;
            justify-content: center;
            color: #000;
            font-size: 14px;
        }
        .buy-btn {
            display: inline-block;
            padding: 12px 30px;
            background-color: #000;
            color: white;
            text-decoration: none;
            border: 2px solid #000;
            border-radius: 4px;
            font-weight: 600;
            font-size: 14px;
            cursor: pointer;
            transition: all 0.3s ease;
        }
        .buy-btn:hover {
            background-color: #333;
            border-color: #333;
        }
        .divider {
            height: 1px;
            background-color: #000;
            margin: 40px 0;
        }
        .footer-section {
            background-color: #f0f0f0;
            padding: 30px 20px;
            text-align: center;
            border-top: 2px solid #000;
        }
        .footer-text {
            font-size: 16px;
            font-weight: 500;
            color: #000;
            margin-bottom: 20px;
        }
        .call-btn {
            display: inline-block;
            padding: 14px 25px;
            background-color: #fff;
            color: #000;
            text-decoration: none;
            border: 2px solid #000;
            border-radius: 4px;
            font-weight: 600;
            font-size: 14px;
            cursor: pointer;
            transition: all 0.3s ease;
        }
        .call-btn:hover {
            background-color: #000;
            color: #fff;
        }
        .call-icon {
            margin-right: 8px;
            font-size: 16px;
        }
        .footer-info {
            font-size: 12px;
            color: #666;
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <!-- Header Banner -->
        <div class="header-banner">
            <h1>🎁 Special Offers Just For You</h1>
        </div>

        <!-- Main Content -->
        <div class="content">
            <div class="greeting">Hello, %s</div>
            
            <div class="intro-text">
                We've curated some amazing products we think you'll love. Discover our latest recommendations tailored especially for you. Don't miss out on these exclusive deals and offers available for a limited time only!
            </div>

            %s

            <div class="divider"></div>

            <!-- Footer Section -->
            <div class="footer-section">
                <div class="footer-text">Have questions or need assistance?</div>
                <a href="http://localhost:8080/call-action?phone=%s" class="call-btn">
                    <span class="call-icon">📞</span>Make a call!
                </a>
                <div class="footer-info">
                    <p>We're here to help Monday - Friday, 9AM - 6PM EST</p>
                </div>
            </div>
        </div>
    </div>
</body>
</html>`

	return fmt.Sprintf(htmlTemplate, req.UserName, productsSection, req.PhoneNumber)
}

// Handler para enviar el correo
func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w) // Habilitar CORS para todas las solicitudes

	// Manejar las solicitudes OPTIONS para CORS
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Asegurarse de que sea un POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar el JSON de la solicitud
	var emailReq EmailRequest
	err := json.NewDecoder(r.Body).Decode(&emailReq)
	if err != nil {
		http.Error(w, "Error al procesar el JSON", http.StatusBadRequest)
		return
	}

	// Construir el contenido del correo
	content := fmt.Sprintf(`
		<h1>%s</h1>
		<h2>%s</h2>
		<p>%s</p>
	`, emailReq.Mail, emailReq.Subject, emailReq.Body)

	// Crear el remitente usando el paquete mail
	sender := mail.NewGmailSender(
		os.Getenv("EMAIL_SENDER_NAME"),
		os.Getenv("EMAIL_SENDER_ADDRESS"),
		os.Getenv("EMAIL_SENDER_PASSWORD"),
	)

	// Enviar el correo
	to := []string{os.Getenv("DESTINATION_EMAIL")}
	attachFiles := []string{} // Puedes agregar archivos si es necesario

	err = sender.SendEmail(emailReq.Subject, content, to, nil, nil, attachFiles)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al enviar el correo: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Correo enviado exitosamente"))
}

// Handler para enviar recomendaciones de productos
func sendRecommendationHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w) // Habilitar CORS para todas las solicitudes

	// Manejar las solicitudes OPTIONS para CORS
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Asegurarse de que sea un POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar el JSON de la solicitud
	var recommendationReq RecommendationRequest
	err := json.NewDecoder(r.Body).Decode(&recommendationReq)
	if err != nil {
		http.Error(w, "Error al procesar el JSON", http.StatusBadRequest)
		return
	}

	// Generar el HTML de las recomendaciones
	htmlContent := generateRecommendationHTML(recommendationReq)

	// Crear el remitente usando el paquete mail
	sender := mail.NewGmailSender(
		os.Getenv("EMAIL_SENDER_NAME"),
		os.Getenv("EMAIL_SENDER_ADDRESS"),
		os.Getenv("EMAIL_SENDER_PASSWORD"),
	)

	// Enviar el correo
	to := []string{recommendationReq.DestinationEmail}
	attachFiles := []string{} // Puedes agregar archivos si es necesario

	err = sender.SendEmail(recommendationReq.Subject, htmlContent, to, nil, nil, attachFiles)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al enviar el correo: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Correo de recomendaciones enviado exitosamente"))
}

// Handler para manejar las llamadas del botón "Make a call"
func callActionHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var phoneNumber string

	if r.Method == http.MethodGet {
		// Obtener el número de teléfono desde los parámetros de la URL
		phoneNumber = r.URL.Query().Get("phone")
		if phoneNumber == "" {
			http.Error(w, "Número de teléfono requerido", http.StatusBadRequest)
			return
		}
	} else if r.Method == http.MethodPost {
		// Decodificar la solicitud JSON para obtener el número de teléfono
		var callReq PhoneCallRequest
		err := json.NewDecoder(r.Body).Decode(&callReq)
		if err != nil {
			http.Error(w, "Error al procesar el JSON", http.StatusBadRequest)
			return
		}
		phoneNumber = callReq.PhoneNumber
	} else {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("Usuario solicitó llamada para el número: %s", phoneNumber)

	// Hacer la llamada a la API externa
	err := makePhoneCall(phoneNumber)
	if err != nil {
		log.Printf("Error al hacer la llamada: %v", err)

		// Responder con HTML para mejor experiencia de usuario desde el email
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		html := fmt.Sprintf(`
			<html>
			<head><title>Error en la llamada</title></head>
			<body style="font-family: Arial, sans-serif; text-align: center; padding: 50px;">
				<h2>❌ Error al procesar la llamada</h2>
				<p>Lo sentimos, ocurrió un error al intentar realizar la llamada al número %s.</p>
				<p>Error: %v</p>
				<a href="javascript:history.back()" style="color: #007bff;">← Volver</a>
			</body>
			</html>
		`, phoneNumber, err)
		w.Write([]byte(html))
		return
	}

	// Responder con HTML para mejor experiencia de usuario desde el email
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	html := fmt.Sprintf(`
		<html>
		<head><title>Llamada iniciada</title></head>
		<body style="font-family: Arial, sans-serif; text-align: center; padding: 50px;">
			<h2>✅ Llamada iniciada exitosamente</h2>
			<p>Se ha iniciado la llamada al número: <strong>%s</strong></p>
			<p>Gracias por usar nuestro servicio.</p>
		</body>
		</html>
	`, phoneNumber)
	w.Write([]byte(html))
}

// Función para hacer la llamada a la API externa
func makePhoneCall(phoneNumber string) error {
	// Preparar el payload para la API externa
	payload := PhoneCallRequest{
		PhoneNumber: phoneNumber,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error al codificar JSON: %v", err)
	}

	// Hacer la petición POST a la API externa
	resp, err := http.Post(
		"http://localhost:8000/api/v1/phonecalls/make_call_body",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("error al hacer petición HTTP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API externa respondió con código: %d", resp.StatusCode)
	}

	log.Printf("Llamada iniciada exitosamente para el número: %s", phoneNumber)
	return nil
}

func main() {
	// Endpoint original para compatibilidad
	http.HandleFunc("/send-email", sendEmailHandler)

	// Nuevo endpoint para recomendaciones de productos
	http.HandleFunc("/recommendations", sendRecommendationHandler)

	// Endpoint para manejar las acciones del botón "Make a call"
	http.HandleFunc("/call-action", callActionHandler)

	fmt.Println("Servidor escuchando en el puerto 8080...")
	fmt.Println("Endpoints disponibles:")
	fmt.Println("  POST /send-email - Envío de correo básico")
	fmt.Println("  POST /recommendations - Envío de recomendaciones de productos")
	fmt.Println("  GET|POST /call-action - Manejo de acciones de llamada")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
