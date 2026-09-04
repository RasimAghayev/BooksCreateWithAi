# Chapter 2 — References

## Book: Разработка приложений в микросервисной архитектуре с нуля

## Primary Sources

1. **Джейсон Андресс. "Защита данных. От авторизации до аудита". 2021.**
   - URL: https://goo.su/7cBj7
   - Context: Auth theory, security threats, password hashing best practices
   - Referenced in: Pages 87-91 (auth theory section), Pages 97-101 (bcrypt discussion)

2. **https://golang-jwt.github.io/jwt/**
   - Title: Go JWT Library Documentation
   - Context: JWT creation, validation, claims parsing in Go
   - Referenced in: Pages 97-101 (Validate, issueTokens methods)

## Secondary Sources (from Chapter 3 context, relevant to Chapter 2 inter-service comms)

3. **https://ru.wikipedia.org/wiki/AMQP**
   - Title: AMQP Protocol (Russian Wikipedia)
   - Context: Message broker protocol referenced for async service communication
   - Referenced in: Chapter 3 (RabbitMQ section)

4. **https://www.rabbitmq.com/maxlength.html**
   - Title: RabbitMQ Queue Length Limits
   - Context: Managing message queue capacity
   - Referenced in: Chapter 3

5. **https://www.confluent.io/blog/kafka-strategy-decision-framework/**
   - Title: Kafka vs RabbitMQ Decision Framework
   - Context: Comparing message brokers
   - Referenced in: Chapter 3
