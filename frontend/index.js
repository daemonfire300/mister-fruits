function SignupRequest(username, password) {
    return {
        username: username,
        password: password,
    }
}

let LoginRequest = SignupRequest

class UserService {
    constructor(url) {
        this.url = url
    }

    async signup(username, password) {
        const resp = await fetch(this.url + '/signup', {
            method: 'POST',
            mode: 'cors',
            cache: 'no-cache',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(SignupRequest(username, password))
        })
        console.log(resp)
        return resp.ok
    }

    async login(username, password) {
        const resp = await fetch(this.url + '/login', {
            method: 'POST',
            mode: 'cors',
            cache: 'no-cache',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(LoginRequest(username, password))
        })
        console.log(resp)
        if (!resp.ok) {
            return false
        }
        let data = await resp.json()
        return data.value // auth token
    }
}

class CartAndProductService {
    constructor(url) {
        this.url = url
    }

    async listProducts(token) {
        const resp = await fetch(this.url + '/products', {
            method: 'GET',
            mode: 'cors',
            cache: 'no-cache',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + token,
            },
        })
        if (resp.ok) {
            return await resp.json()
        }
        return []
    }
}

var userSvc = new UserService('http://localhost:8181')
var cartAndProductSvc = new CartAndProductService('http://localhost:8181')

var app = new Vue({
    el: '#shop',
    data: {
        content: 'wooooop',
        warning: '',
        username: '',
        password: '',
        authToken: '',
        userSvc: userSvc,
        cartAndProductSvc: cartAndProductSvc,
        products: [],
    },
    mounted() {
        if (localStorage.authToken) {
            this.authToken = localStorage.authToken;
            this.loadProducts()
        }
    },
    methods: {
        submitLogin: async function (formEvent) {
            formEvent.preventDefault()
            if (this.password == "" || this.username == "") {
                this.warning = "username and password not set"
                console.log("username and password not set")
            } else {
                let login = await this.userSvc.login(this.username, this.password)
                if (!login) {
                    this.warning = "Login failed"
                    console.log("login failed")
                    return
                }
                console.log("login successful")
                this.authToken = login
                localStorage.authToken = login
                this.loadProducts()
            }

        },
        submitSignup: async function (formEvent) {
            formEvent.preventDefault()
            if (this.password == "" || this.username == "") {
                this.warning = "username and password not set"
                console.log("username and password not set")
            } else {
                let signupRes = await this.userSvc.signup(this.username, this.password)
                if (!signupRes) {
                    this.warning = "Signup failed"
                    return
                }
            }
        },
        loadProducts: async function () {
            this.products = await this.cartAndProductSvc.listProducts(this.authToken)
        }
    }
})