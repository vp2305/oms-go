# oms-go

## Tools
- `cosmtrek/air`: Live reload for Go apps
- `gRPC` for communication between services
- `RabbitMQ` as message broker
    - Retry mechanism
- `Docker` with docker compose
- `MongoDB` as storage layer
- `Jaeger` for service tracing
- `HashiCorp's Consul` for service discovery
- `Stripe` for payments

## Architecture
![alt text](images/architecture.png)

## Services Breakdown
Order Service: 
- Validate order details -> talk with stock service
- CRUD of orders
- Initiates the Payment Flow -> by sending an event

Stock Service:
- Handles stock
- Validate order quantities
- (might return items as menu)

Menu Service:
- stores items as menu

Payment Service:
- Initiates a payment with a 3rd party provider
- Produces an order Paid/Cancelled event to (orders, stock and kitchen)

Kitchen Service
- Long running process of a "simulated kitchen staff"


## Notes
### Monolithic Architecture
Pros:
- Monolith architecture is more simpler which is a good thing! Specifically to deploy compared to distributed systems since we don't have to split the application into multiple parts that coordinate with each other.
- Making the codebase smaller and easier to develop and maintain. You have all the code in one repository project making it easier to understand the system as a whole and make changes across the application.
Cons:
- At a certain point the application grows so much that it becomes a bottleneck to deploy.
- Scalability: Monolithic architectures can be challenging to scale, as the entire application needs to be replicated to handle increased load.
- This is especially the case as the engineering team increases and the deploys are more frequent.

### Microservices Architecture
Pros:
- Scalability: Microservices enable better horizontal scaling, as you can scale individual services as needed.
- Fault tolerance: If one service fails, it doesn’t necessarily mean the entire application will fail.
- Smaller services: Smaller services are easier to understand, develop, and maintain.
- Decoupling: Microservices are loosely coupled, so you can update one service without affecting the others.
- Technology diversity: You can use different technologies for different services.
Cons:
- It's more resource intensive because when a system consists of multiple components, they are divided among multiple servers. They need to communicate and that communication increases the latency.
- Harder to debug
- Consistency: data is often scattered among the 
- Possible duplication of code and features
- Observability: Which because the system is made of the parts we need to collect more data, including traces, logs and other metrics.