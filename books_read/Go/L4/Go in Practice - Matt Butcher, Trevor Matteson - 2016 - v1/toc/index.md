# Go in Practice — Mündəricat (TOC)

**Müəlliflər:** Matt Butcher, Trevor Matteson
**Nəşriyyat:** Manning Publications Co., 2016 · **ISBN:** 9781633430075
**Səhifə:** 314 (PDF) · **Book pages:** 280 + indeks
**Struktur:** 4 hissə · 11 chapter · **70 texnika (TECHNIQUE)**

> PDF offset: book page N → PDF page N+23

## Hissələr

| Part | Ad | Mövzu | Chapter |
|---|---|---|---|
| 1 | Background and Fundamentals | CLI, konfiq, server, concurrency | 1-3 |
| 2 | Well-Rounded Applications | xətalar, debug/test | 4-5 |
| 3 | An Interface for Your Applications | template-lər, formalar, web services | 6-8 |
| 4 | Taking Your Applications to the Cloud | cloud, microservices, reflection | 9-11 |

## Chapter-lər

| Ch | Başlıq | Book səh. | PDF səh. | Texnikalar |
|----|--------|-----------|----------|------------|
| 1 | Getting into Go | 3-25 | 26-48 | — |
| 2 | A solid foundation | 27-58 | 50-81 | 1-9 |
| 3 | Concurrency in Go | 59-83 | 82-106 | 10-15 |
| 4 | Handling errors and panics | 87-111 | 110-134 | 16-21 |
| 5 | Debugging and testing | 113-143 | 136-166 | 22-31 |
| 6 | HTML and email template patterns | 147-166 | 170-189 | 32-38 |
| 7 | Serving and receiving assets and forms | 168-193 | 191-212 | 39-48 |
| 8 | Working with web services | 194-213 | 217-236 | 49-55 |
| 9 | Using the cloud | 217-234 | 240-257 | 56-61 |
| 10 | Communication between cloud services | 235-252 | 258-275 | 62-65 |
| 11 | Reflection and code generation | 253-280 | 276-303 | 66-70 |

## 70 Texnikanın xəritəsi

### Ch 2 — A solid foundation (1-9)
1. GNU/UNIX-style command-line arguments · 2. Avoiding CLI boilerplate (frameworks) · 3. Using configuration files · 4. Configuration via environment variables · 5. Graceful shutdowns using manners · 6. Matching paths to content · 7. Handling complex paths with wildcards · 8. URL pattern matching · 9. Faster routing (without the work)

### Ch 3 — Concurrency in Go (10-15)
10. Using goroutine closures · 11. Waiting for goroutines (WaitGroup) · 12. Locking with a mutex · 13. Using multiple channels · 14. Closing channels · 15. Locking with buffered channels

### Ch 4 — Handling errors and panics (16-21)
16. Minimize the nils · 17. Custom error types · 18. Error variables · 19. Issuing panics · 20. Recovering from panics · 21. Trapping panics on goroutines

### Ch 5 — Debugging and testing (22-31)
22. Logging to an arbitrary writer · 23. Logging to a network resource · 24. Handling back pressure in network logging · 25. Logging to the syslog · 26. Capturing stack traces · 27. Using interfaces for mocking or stubbing · 28. Verifying interfaces with canary tests · 29. Benchmarking Go code · 30. Parallel benchmarks · 31. Detecting race conditions

### Ch 6 — HTML and email template patterns (32-38)
32. Extending templates with functions · 33. Caching parsed templates · 34. Handling template execution failures · 35. Nested templates · 36. Template inheritance · 37. Mapping data types to templates · 38. Generating email from templates

### Ch 7 — Serving and receiving assets and forms (39-48)
39. Serving subdirectories · 40. File server with custom error pages · 41. Caching file server · 42. Embedding files in a binary · 43. Serving from an alternative location · 44. Accessing multiple values for a form field · 45. Uploading a single file · 46. Uploading multiple files · 47. Verify uploaded file is allowed type · 48. Incrementally saving a file

### Ch 8 — Working with web services (49-55)
49. Detecting timeouts · 50. Timing out and resuming with HTTP · 51. Custom HTTP error passing · 52. Reading custom errors · 53. Parsing JSON without knowing the schema · 54. API version in the URL · 55. API version in content type

### Ch 9 — Using the cloud (56-61)
56. Working with multiple cloud providers · 57. Cleanly handling cloud provider errors · 58. Gathering information on the host · 59. Detecting dependencies · 60. Cross-compiling · 61. Monitoring the Go runtime

### Ch 10 — Communication between cloud services (62-65)
62. Reusing connections · 63. Faster JSON marshal and unmarshal · 64. Using protocol buffers · 65. Communicating over RPC with protocol buffers

### Ch 11 — Reflection and code generation (66-70)
66. Switching based on type and kind · 67. Discovering whether a value implements an interface · 68. Accessing fields on a struct · 69. Processing tags on a struct · 70. Generating code with go generate

## Layihə davamlılığı (project_continuity)

```json
{
  "project_continuity": false,
  "note": "Kitab texnika-kolleksiyasıdır: hər TECHNIQUE müstəqil, kiçik, özünü sübut edən nümunədir. Kumulyativ layihə yoxdur — project-state/ lazım deyil. 70 texnika çərçivəsində fərqli pattern-lər müstəqil göstərilir."
}
```
