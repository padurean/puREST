# github.com/padurean/purest

Welcome to the puREST generated docs!

## Routes

<details>
<summary>`/`</summary>

- [RequestID]()
- [RealIP]()
- [Recoverer]()
- [Heartbeat.func1]()
- [Timeout.func1]()
- [WithValue.func1]()
- [NewHandler.func1]()
- [RemoteAddrHandler.func1]()
- [UserAgentHandler.func1]()
- [RefererHandler.func1]()
- [RequestIDHandler.func1]()
- [WithValue.func1]()
- **/**
	- _GET_
		- [Router.setupRoutes.func1](/server/router.go#L106)

</details>
<details>
<summary>`/api/*/v1/*/users/*`</summary>

- [RequestID]()
- [RealIP]()
- [Recoverer]()
- [Heartbeat.func1]()
- [Timeout.func1]()
- [WithValue.func1]()
- [NewHandler.func1]()
- [RemoteAddrHandler.func1]()
- [UserAgentHandler.func1]()
- [RefererHandler.func1]()
- [RequestIDHandler.func1]()
- [WithValue.func1]()
- **/api/***
	- **/v1/***
		- **/users/***
			- **/**
				- _POST_
					- [UserCreate](/controllers/user.go#L149)
				- _GET_
					- [paginate](/server/router.go#L80)
					- [UserList](/controllers/user.go#L201)

</details>
<details>
<summary>`/api/*/v1/*/users/*/sign-in/{usernameOrEmail}`</summary>

- [RequestID]()
- [RealIP]()
- [Recoverer]()
- [Heartbeat.func1]()
- [Timeout.func1]()
- [WithValue.func1]()
- [NewHandler.func1]()
- [RemoteAddrHandler.func1]()
- [UserAgentHandler.func1]()
- [RefererHandler.func1]()
- [RequestIDHandler.func1]()
- [WithValue.func1]()
- **/api/***
	- **/v1/***
		- **/users/***
			- **/sign-in/{usernameOrEmail}**
				- _POST_
					- [UserCtx](/controllers/user.go#L86)
					- [UserSignIn](/controllers/user.go#L182)

</details>
<details>
<summary>`/api/*/v1/*/users/*/{id}/*`</summary>

- [RequestID]()
- [RealIP]()
- [Recoverer]()
- [Heartbeat.func1]()
- [Timeout.func1]()
- [WithValue.func1]()
- [NewHandler.func1]()
- [RemoteAddrHandler.func1]()
- [UserAgentHandler.func1]()
- [RefererHandler.func1]()
- [RequestIDHandler.func1]()
- [WithValue.func1]()
- **/api/***
	- **/v1/***
		- **/users/***
			- **/{id}/***
				- [UserCtx](/controllers/user.go#L86)
				- **/**
					- _DELETE_
						- [UserDelete](/controllers/user.go#L267)
					- _GET_
						- [UserGet](/controllers/user.go#L260)
					- _PUT_
						- [UserUpdate](/controllers/user.go#L226)

</details>
<details>
<summary>`/swagger/*`</summary>

- [RequestID]()
- [RealIP]()
- [Recoverer]()
- [Heartbeat.func1]()
- [Timeout.func1]()
- [WithValue.func1]()
- [NewHandler.func1]()
- [RemoteAddrHandler.func1]()
- [UserAgentHandler.func1]()
- [RefererHandler.func1]()
- [RequestIDHandler.func1]()
- [WithValue.func1]()
- **/swagger/***
	- _GET_
		- [github.com/swaggo/http-swagger.Handler.func1]()

</details>

Total # of routes: 5
