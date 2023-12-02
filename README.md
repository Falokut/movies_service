# Content

+ [Configuration](#configuration)
+ [Metrics](#metrics)
+ [Docs](#docs)
+ [Author](#author)
+ [License](#license)
---------

# Configuration

1. [Configure movies_db](movies_db/README.md#Configuration)
2. Create .env on project root dir  
Example env:
```env
REDIS_PASSWORD=redispass
REDIS_AOF_ENABLED=no
DB_PASSWORD=Password
```


# Metrics
The service uses Prometheus and Jaeger and supports distribution tracing

# Docs
[Swagger docs](swagger/docs/movies_service_v1.swagger.json)

# Author

- [@Falokut](https://github.com/Falokut) - Primary author of the project

# License

This project is licensed under the terms of the [MIT License](https://opensource.org/licenses/MIT).

---