# Payment Service

The **Payment Service** is a core component of the PharmaKart platform, responsible for handling order payments securely using Stripe. It manages payment processing, refunds, and transaction status updates, ensuring seamless integration with the order system.

---

## Table of Contents
1. [Overview](#overview)
2. [Features](#features)
3. [Prerequisites](#prerequisites)
4. [Setup and Installation](#setup-and-installation)
5. [Running the Service](#running-the-service)
6. [Environment Variables](#environment-variables)
7. [Contributing](#contributing)
8. [License](#license)

---

## Overview

The Payment Service handles:
- Payment processing using **Stripe**.
- Payment status retrieval.
- Refund management.
- Secure transaction handling with Stripe API integration.

It is built using **gRPC** for communication and **Stripe API** for payment handling.

---

## Features

- **Payment Processing**:
  - Create a payment intent using Stripe.
  - Handle successful and failed transactions.
- **Refund Management**:
  - Initiate and process refunds for orders.
- **Payment Status Retrieval**:
  - Fetch payment details and status updates.

---

## Prerequisites

Before setting up the service, ensure you have the following installed:
- **Docker**
- **Go** (for building and running the service)
- **Protobuf Compiler** (`protoc`) for generating gRPC/protobuf files
- **Stripe Account** (for API keys and sandbox testing)
- **Stripe CLI** (optional, for local testing)

---

## Setup and Installation

### 1. Clone the Repository
Clone the repository and navigate to the payment service directory:
```bash
git clone https://github.com/PharmaKart/payment-svc.git
cd payment-svc
```

### 2. Generate Protobuf Files
Generate the protobuf files using the provided `Makefile`:
```bash
make proto
```

### 3. Install Dependencies
Run the following command to ensure all dependencies are installed:
```bash
go mod tidy
```

### 4. Build the Service
To build the service, run:
```bash
make build
```

### 5. Run local stripe forwarder
```bash
stripe login
stripe listen --forward-to localhost:8080/payment/webhook
```

---

## Running the Service

### Option 1: Run Using Docker
To run the service using Docker, execute:
```bash
docker run -p 50054:50054 pharmakart/payment-svc
```

### Option 2: Run Using Makefile
To run the service directly using Go, execute:
```bash
make run
```

The service will be available at:
- **gRPC**: `localhost:50054`

---

## Environment Variables

The service requires the following environment variables. Create a `.env` file in the `payment-svc` directory with the following:

```env
ORDER_DB_HOST=postgres
ORDER_DB_PORT=5432
ORDER_DB_USER=postgres
ORDER_DB_PASSWORD=yourpassword
ORDER_DB_NAME=pharmakartdb
ORDER_SERVICE_ADDR=localhost:50053
STRIPE_SECRET_KEY=your-stripe-secret-key
GATEWAY_URL=http://localhost:8080
```

---

## Contributing

Contributions are welcome! Please follow these steps:
1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Submit a pull request with a detailed description of your changes.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Support

For any questions or issues, please open an issue in the repository or contact the maintainers.

