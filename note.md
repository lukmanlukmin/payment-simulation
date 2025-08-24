# Notes & Improvements

## Assumptions
- **3rd-Party Success Rate**: The external provider call is simulated with a 98% success probability.  
  - This means ~2% of transactions are expected to fail, which triggers refund and retry logic.  
- **Retry Logic**: Optimistic transaction retries are limited to `maxRetries = 5` with exponential backoff.  
  - Assumes that conflicts or transient DB errors resolve within these retries.  
- **Time Constraints / Delays**:  
  - Sleep/delay simulates provider latency (300–1500ms). Assumes production provider has similar response time distribution.  
- **Minimal Testing Scope**:  
  - Only core service logic is unit-tested; assumes repository, Kafka, and HTTP layers behave as expected based on mocks.  
- **Load Testing Baseline**:  
  - Simulated >10 RPS with k6; assumes performance metrics are representative of light to moderate load.  
- **Environment Assumptions**:  
  - PostgreSQL and Kafka connections succeed, no network interruptions.  
  - Decimal arithmetic does not overflow for tested transaction amounts.

## Unit Testing Coverage
- Only minimal unit tests implemented for core service logic.
- Testing focuses mainly on:
  - `SubmitTransfer` core flow
  - `ProcessTransaction` basic success / status checks
- **Reason for limited coverage**:
  - Time constraint
  - Some parts depend on external systems (Kafka, DB transactions) which were mocked
  - Complex retry/backoff and optimistic transaction logic are partially tested

## Optimistic Transaction Handling
- The service uses **optimistic transaction** pattern:
  - Deduct/Credit balance only if `version` matches
  - Retry on conflicts (`maxRetries` with exponential backoff)
- This ensures atomicity without database locks but makes full unit testing tricky, especially for concurrency and conflict scenarios.

## Points for Future Improvement
- **Repository Layer**: more extensive unit tests for DB queries and edge cases
- **Message Broker Integration**: full test coverage for Kafka publish/consume
- **HTTP Layer**: input validation, error handling, and integration tests
- **Retry / Conflict Handling**: simulate concurrent transactions and optimistic lock conflicts
- **Logging & Tracing**: 
  - Implement structured logging and distributed tracing (e.g., OpenTelemetry) for observability.
  - Helps to monitor performance metrics such as response time, success rate, and error rate under load.
  - Would make debugging and analyzing retries, optimistic transaction conflicts, and 3rd-party failures much easier.
- **Error Handling / Alerting**: improve propagation of errors for failed transactions and monitoring

## Summary
- Core functionality tested minimally; optimistic transaction flow is partially covered
- With more time, test coverage and observability would be expanded for production-readiness
