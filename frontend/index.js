function uuidv4() {
    return ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, c =>
        (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
    );
}

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

    async createCart(token, id) {
        const resp = await fetch(this.url + '/carts/' + id, {
            method: 'POST',
            mode: 'cors',
            cache: 'no-cache',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + token,
            },
            body: JSON.stringify({
                id: id,
                items: [],
                coupons: [],
            })
        })
        return resp.ok
    }

    async fetchCart(token, id) {
        const resp = await fetch(this.url + '/carts/' + id, {
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
    }

    async updateCart(token, id, add, remove) {
        const resp = await fetch(this.url + '/carts/' + id, {
            method: 'PATCH',
            mode: 'cors',
            cache: 'no-cache',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + token,
            },
            body: JSON.stringify({
                'add': add,
                'remove': remove,
            })
        })
        return resp.ok
    }

    async checkoutCart(token, id) {
        const resp = await fetch(this.url + '/carts/' + id + '/checkout', {
            method: 'POST',
            mode: 'cors',
            cache: 'no-cache',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + token,
            }
        })
        return resp.ok
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
        cartID: '',
        cart: null,
        currentCheckoutSum: 0.00
    },
    mounted() {
        if (localStorage.authToken) {
            this.authToken = localStorage.authToken;
            this.loadProducts()
            if (localStorage.cartID) {
                this.cartID = localStorage.cartID
                if (this.cartID !== '') {
                    this.loadCart()
                }
            }
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
        },
        loadCart: async function () {
            this.cart = await this.cartAndProductSvc.fetchCart(this.authToken, this.cartID)
            if (this.cart != null) {
                let vueRoot = this
                vueRoot.currentCheckoutSum = 0.00
                this.cart.sets.forEach((set) => {
                    vueRoot.currentCheckoutSum += Number(set.grossPrice)
                })
            }
        },
        addOrRemoveProduct: async function (productID, addFlag) {
            if (this.cart == null) { // create cart if it does not exist
                let id = uuidv4()
                let res = await this.cartAndProductSvc.createCart(this.authToken, id)
                if (!res) {
                    console.log("error could not create cart...")
                    return
                }
                this.cartID = id
                localStorage.cartID = this.cartID
            }
            let remove = [];
            let add = [];
            let p = { product: { id: productID }, quantity: 1 };
            if (addFlag) {
                add.push(p)
            } else {
                remove.push(p)
            }
            let res = await this.cartAndProductSvc.updateCart(this.authToken, this.cartID, add, remove)
            if (res) {
                this.loadCart()
            }
        },
        checkout: async function () {
            res = await this.cartAndProductSvc.checkoutCart(this.authToken, this.cartID)
            if (res) {
                localStorage.cartID = ''
                this.cartID = ''
                this.cart = null
            }
        },
        logout: function () {
            localStorage.authToken = ''
            this.authToken = ''
            localStorage.cartID = ''
            this.cartID = ''
            this.cart = null
        }
    }
})