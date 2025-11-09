# CoAP Theos DNS

This project is a CoAP-based DNS resolver that converts CoAP requests to DNS queries and returns the resolved IP address as a response. It is designed for lightweight communication using the CoAP protocol.

## ENVs

- `LOOKUP_SERVER`: Specifies the DNS server to resolve requests. Defaults to `8.8.8.8:53` if not set.

## PORT

- `5688`: The server listens on this port for incoming CoAP requests.

## API Path

| **Method** | **Path** | **Description**                              |
|------------|----------|----------------------------------------------|
| `POST`     | `/ip`    | Resolves a domain name to its IPv4 address. |

## Usage

1. **Set the DNS server (optional):**  
   You can specify the DNS server using the `LOOKUP_SERVER` environment variable. If not set, the default is `8.8.8.8:53`.

   ```bash
   export LOOKUP_SERVER="1.1.1.1:53"
   ```

2. **Run the server:**  
   Build and run the server using the provided `Dockerfile` or directly with Go.

   Example using Docker:
   ```bash
   docker build -t coap-dns-resolver .
   docker run -p 5688:5688 -e LOOKUP_SERVER="1.1.1.1:53" coap-dns-resolver
   ```

3. **Send a CoAP request:**  
   Use a CoAP client to send a request to the `/ip` endpoint with the domain name as the payload.

   Example using `coap-client`:
   ```bash
   coap-client-notls -m post -e "example.com" coap://localhost:5688/ip
   ```

   The server will respond with the resolved IP address or `NXDOMAIN` if the lookup fails.
