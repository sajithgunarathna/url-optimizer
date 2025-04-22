# URL Optimizer

Web Page Analyzer is a web application that analyzes web pages based on user-provided URLs. It extracts information such as the HTML version, title, headings count, link details, and login forms.

The backend is built using Golang (Gin Framework)


## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/your-username/url-optimizer.git
    ```
2. Navigate to the project directory:
    ```bash
    cd url-optimizer
    ```
3. Install dependencies:
    ```bash
    go mod tidy
    ```

## Usage

1. Start the development server:
    ```bash
    cd cmd/web-analyzer
    go run main.go
    ```
2. Access the service at `http://localhost:8080`.

## API Endpoints

- **POST /analyze**  
  Analyze a given web page 

- **GET /status**  
  Check service health status

- **GET /urls**  
  Get analyzed URLs history

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
