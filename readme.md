# DoitPay Take Home Test
This project is a **Payment Service** built with **Golang**, using:  
- **PostgreSQL** as the database  
- **Kafka** as the message broker  
- **Docker Compose** for infrastructure orchestration  

## Infrastructure

All services are containerized using Docker Compose and run on a shared network. The configuration is defined in the **docker-compose.yml** file.

The project is already preconfigured but you can modify the configuration according to your needs.

- The API service is preconfigured to listen on port `:80`.
- Service configuration file location: `src/config.yaml`
- You can modify any configuration (e.g., database URL, Kafka broker, ports) according to your requirements.

## Services

The project contains two main services:
1. **API Service** → Handles incoming HTTP requests.
2. **Worker Service** → Processes background jobs.

## Usage

### Prerequisites
Make sure you have installed:
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Running the Project without Tests
To start the project and run the load tests:

```bash
docker-compose up
```
### Stopping the Project
To stop and clean up all containers and volumes:

```bash
docker compose down -v
```

### Running the Project with Tests
To start the project and run the load tests:

```bash
sh run-app.sh
```

This will:
1. Start all required services (PostgreSQL, Kafka, Zookeeper, API, Worker, and K6).
2. Wait until services are ready.
3. Run load tests automatically with **k6**.

### Stopping the Project
To stop and clean up all containers and volumes:

```bash
docker compose down -v
```

This ensures you start from a clean state when running tests again.

## Checking Results

You can verify the test results directly by checking the database, which is configured inside the `docker-compose.yml` file.  
The database will contain the updated state after the load tests run.

## Load Test Details

The load test is powered by [k6](https://k6.io/).  
The script is located in `/scripts/transaction.js` and simulates transfer requests to:

```
POST /transfer
```

### Example Payload
```json
{
  "amount": 10000,
  "bank_code": "BCA",
  "beneficiary_account": "12345",
  "beneficiary_name": "Test Name",
  "note": "test transfer"
}
```

### Possible Response Status Codes
- `202` → Accepted  
- `400` → Bad Request  
- `500` → Internal Server Error  

K6 will summarize the number of responses per status code in the output.

---

## Logs

To view logs in real-time:

```bash
docker compose logs -f doitpay-api doitpay-worker
```

Or for all services:

```bash
docker compose logs -f
```

---

## Recap
- Run project only → `docker compose up`  
- Run project with test → `sh run-app.sh`  
- Rerun test → `docker compose down -v && sh run-app.sh`  
- Check logs → `docker compose logs -f`  

---