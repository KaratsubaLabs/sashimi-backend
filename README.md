# サシミ - Sashimi Backend
Backend for からつばLABS' **project sashimi** - a status monitor for all our services

## RUNNING FOR DEVELOPMENT
```
go run sashimi.go
```

## TODO

- [ ] Clean up structs (Seperate database models and api schemas, add "getURL()" db function)
- [ ] Add Interactive CLI (Add/delete service, run/stop sashimi, sashimi should autostart)
- [ ] Review/ensure consistent error handling
- [ ] Dockerize