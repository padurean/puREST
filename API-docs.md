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
		- [Router.setupRoutes.func1]()

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
					- [authenticate]()
					- [UserCreate]()
				- _GET_
					- [authenticate]()
					- [paginate]()
					- [UserList]()

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
					- [UserCtx]()
					- [UserSignIn]()

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
				- [UserCtx]()
				- **/**
					- _DELETE_
						- [UserDelete]()
					- _GET_
						- [UserGet]()
					- _PUT_
						- [UserUpdate]()

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
