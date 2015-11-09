GIF backend
===============

*Frontend: https://github.com/Vrenc/gif-frontend*

### Configuration
Create a config file named `config.json`. See `config.example.json` for an example. `frontend_host` is used for the CORS header, use `*` if you're lazy. Make sure the front and backend HMAC keys match.

### Docker Compose example
```
frontend:
    image: gif-frontend
    ports:
        - "80:80"
backend:
    image: gif-backend
    ports:
        - "3000:3000"
    links:
        - redis
redis:
    image: redis
```
