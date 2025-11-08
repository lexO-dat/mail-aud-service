# Email Sender API

This project provides a REST API for sending emails using SMTP, with a focus on sending emails through Gmail. It includes functionality for basic email sending and product recommendations with beautiful HTML templates.

## Getting Started

To use this package, you'll need to obtain an application-specific password for the Gmail account you want to send emails from. This password is required for authentication when sending emails.

First init the mod:
```bash
go mod init [project name] && go mod tidy
```

Then you can get the two go modules from github:

```bash
go get github.com/jordan-wright/email && go get github.com/stretchr/testify/require
```

### Obtaining Application-Specific Password

To obtain an application-specific password for Gmail:

1. Go to your sender gmail Account settings: [https://myaccount.google.com/](https://myaccount.google.com/)
2. Click on "Security" in the left sidebar.
3. You have to activate the two way factor
4. Go to the two way factor menu and search the application passwords
5. Then you have to create an application and copy the code, that is your password


### Usage

1. Create a `GmailSender` instance using `NewGmailSender` function, providing your name, Gmail address, and the application-specific password.
2. Call the `SendEmail` method of the `GmailSender` instance to send emails.

Here is an example:

```go
package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const EMAIL_SENDER_NAME = "Example Name"
const EMAIL_SENDER_ADRESS = "examplemail@gmail.com"
const EMAIL_SENDER_PASSWORD = "password that google give you"

func TestNewGmailSender(t *testing.T) {

	sender := NewGmailSender(EMAIL_SENDER_NAME, EMAIL_SENDER_ADRESS, EMAIL_SENDER_PASSWORD)

	subject := "Test subject"

	content := `
	<h1>EXAMPLE</h1>
	<p>Hello world!</p>
	`

	to := []string{"exampledestination@gmail.com"}

	attachFiles := []string{"rute to your file to send"}

	err := sender.SendEmail(subject, content, to, nil, nil, attachFiles)

	require.NoError(t, err)

}

```

## API Endpoints

### 1. Basic Email Sending

**Endpoint:** `POST /send-email`

Send a basic email with simple content.

**Request Body:**
```json
{
  "mail": "Sender Name",
  "subject": "Email Subject",
  "body": "Email content"
}
```

### 2. Product Recommendations

**Endpoint:** `POST /recommendations`

Send a beautifully formatted email with product recommendations using a professional HTML template.

**Request Body:**
```json
{
  "user_name": "Juan",
  "subject": "Productos especiales seleccionados para ti",
  "products": [
    {
      "name": "Auriculares Premium Bluetooth",
      "description": "Experimenta la excelencia con nuestros auriculares mejor valorados. Perfectos para uso diario con calidad excepcional y valor inigualable.",
      "image": "https://ejemplo.com/auriculares.jpg",
      "buy_url": "https://tienda.com/auriculares-premium"
    },
    {
      "name": "Smartwatch Deportivo",
      "description": "Eleva tu estilo de vida con este artículo más vendido. Confiado por miles de clientes satisfechos en todo el mundo.",
      "image": "https://ejemplo.com/smartwatch.jpg",
      "buy_url": "https://tienda.com/smartwatch-deportivo"
    }
  ],
  "call_to_action_url": "https://api.ejemplo.com/contact",
  "phone_number": "+56973756474",
  "destination_email": "cliente@ejemplo.com"
}
```

**Features:**
- Professional HTML email template
- Product showcase with images
- Call-to-action button that initiates phone calls
- Responsive design
- Modern styling
- Integration with external phone call API

### 3. Phone Call Action

**Endpoint:** `GET /call-action?phone={phone_number}` or `POST /call-action`

Handle phone call requests from the email "Make a call" button. This endpoint makes a call to an external API to initiate phone calls.

**GET Request:**
```
GET /call-action?phone=+56973756474
```

**POST Request Body:**
```json
{
  "phone_number": "+56973756474"
}
```

**External API Integration:**
The service automatically calls: `http://localhost:8000/api/v1/phonecalls/make_call_body`

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
EMAIL_SENDER_NAME=Tu Nombre
EMAIL_SENDER_ADDRESS=tu-email@gmail.com
EMAIL_SENDER_PASSWORD=tu-contraseña-de-aplicación
DESTINATION_EMAIL=destino@ejemplo.com
```

## Running the Server

```bash
go run main.go
```

The server will start on port 8080 with the following endpoints:
- `POST /send-email` - Basic email sending
- `POST /recommendations` - Product recommendations email
- `GET|POST /call-action` - Phone call initiation

## Testing with cURL

### Basic Email:
```bash
curl -X POST http://localhost:8080/send-email \
  -H "Content-Type: application/json" \
  -d '{
    "mail": "Test User",
    "subject": "Test Subject",
    "body": "Test message content"
  }'
```

### Product Recommendations:
```bash
curl -X POST http://localhost:8080/recommendations \
  -H "Content-Type: application/json" \
  -d @example-recommendation-request.json
```

### Phone Call Action:
```bash
# Using GET request
curl "http://localhost:8080/call-action?phone=%2B56973756474"

# Using POST request
curl -X POST http://localhost:8080/call-action \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "+56973756474"
  }'
```
