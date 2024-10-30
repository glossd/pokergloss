# Auth
Authentication middleware for go microservices

### Installing

Create `~/.netrc` file with content
```
machine machine gitlab.com
  login your-username@gitlab.com
  password your-personal-access-token
```

Now you can download this library
```
go get --insecure gitlab.com/pokerblow/go-auth
```

### Usage
You need to invoke `auth.Init()` to initialize authentication.
To set user Identity set Middleware on your gin router.  
To get user Identity use `authid.Id`.

### Options
To disable jwt verification export environment variable `PG_JWT_VERIFICATION_DISABLE=true`