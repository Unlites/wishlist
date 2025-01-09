// Import necessary modules
import { createApp } from 'vue';
import axios from 'axios';
import { createRouter, createWebHistory } from 'vue-router';
import '../assets/styles.css';

// Define API base URL
const API_BASE_URL = 'http://127.0.0.1:8080/api/v1';

// Helper for Axios instance with interceptor
const api = axios.create();
api.interceptors.response.use(
    response => response,
    error => {
        $('.alert').alert();
        if (error.response && error.response.status === 401) {
            router.push('/auth');
        }
        return Promise.reject(error);
    }
);

// Auth Page Component
const AuthPage = {
    data() {
        return {
            user: {
                name: '',
                password: '',
                token: ''
            },
            error: null
        };
    },
    methods: {
        async register() {
            try {
                await axios.post(`${API_BASE_URL}/users/register`, {
                    name: this.user.name,
                    password: this.user.password
                });
                this.error = null;
            } catch (error) {
                if (error.response && error.response.data) {
                    if (error.response.status === 400) {
                        alert('Ошибка валидации: логин и пароль должны быть не менее 4 символов либо такой пользователь уже существует');
                    } else {
                        alert('Ошибка регистрации: ' + error.response.data.message);
                    }
                } else {
                    alert('Ошибка регистрации, обратитесь к администратору');
                }
            }
        },
        async login() {
            try {
                const response = await axios.post(`${API_BASE_URL}/users/login`, {
                    name: this.user.name,
                    password: this.user.password
                });
                this.user.token = response.data.token;
                const tokenPayload = JSON.parse(atob(this.user.token.split('.')[1]));                
                localStorage.setItem('token', this.user.token);
                localStorage.setItem('user_id', tokenPayload.sub);
                this.$router.push('/');
            } catch (error) {
                if (error.response && error.response.data) {
                    if (error.response.status === 401) {
                        alert('Ошибка аутентификации: неверный логин или пароль');
                    } else {
                        alert('Ошибка аутентификации: ' + error.response.data.message);
                    }
                } else {
                    alert('Ошибка аутентификации, обратитесь к администратору');
                }
            }  
        }
    },
    template: `
    <div class="auth-page">
        <h1>Authentication</h1>
        <div v-if="error" class="error">{{ error }}</div>
        <input v-model="user.name" placeholder="Name" />
        <input v-model="user.password" type="password" placeholder="Password" />
        <button @click="register">Register</button>
        <button @click="login">Login</button>
    </div>
    `
};

// Main App Component
const App = {
    data() {
        return {
            wishes: [],
            newWish: {
                title: '',
                description: ''
            },
            updateWish: {
                id: null,
                title: '',
                description: '',
                is_reserved: false
            },
            error: null
        };
    },
    methods: {
        async fetchWishes() {
            try {
                const response = await api.get(`${API_BASE_URL}/users/${this.$route.params.user_id}/wishes`, {
                    headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
                });
                this.wishes = response.data;
            } catch (error) {
                alert('Ошибка загрузки списка желаний');
            }
        },
        async addWish() {
            try {
                const response = await api.post(`${API_BASE_URL}/users/${this.$route.params.user_id}/wishes`, this.newWish, {
                    headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
                });
                this.wishes.push(response.data);
                this.newWish = { title: '', description: '' };
                await this.fetchWishes();
            } catch (error) {
                alert('Ошибка добавления желания');
            }
        },
        async updateWishDetails() {
            try {
                await api.put(`${API_BASE_URL}/users/${this.$route.params.user_id}/wishes/${this.updateWish.id}`, this.updateWish, {
                    headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
                });
                this.updateWish = { id: null, title: '', description: '', is_reserved: false };
                await this.fetchWishes();
            } catch (error) {
                alert('Ошибка обновления желания'); 
            }
        },
        async deleteWish(wishId) {
            try {
                await api.delete(`${API_BASE_URL}/users/${this.$route.params.user_id}/wishes/${wishId}`, {
                    headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
                });
                await this.fetchWishes();
            } catch (error) {
                alert('Ошибка удаления желания');
            }
        },
        logout() {
            localStorage.removeItem('token');
            localStorage.removeItem('user_id');
            this.$router.push('/auth');
        },
        isOwnUser() {
            return localStorage.getItem('user_id') == this.$route.params.user_id;
        }
    },
    mounted() {
        this.fetchWishes();
    },
    template: `
    <div class="app p-3">
        <h1 class="text-center">Wishlist</h1>
        <hr>

        <div>
            <div v-if="isOwnUser()" class="text-center m-5 row col-6 col-lg-3 mx-auto border border-2 border-success px-3 py-5 rounded">
                <h3>Новое желание</h3>
                <input class="m-1" v-model="newWish.title" placeholder="Название" />
                <input class="m-1" v-model="newWish.description" placeholder="Ссылка/описание" />
                <button @click="addWish" class="m-1 btn btn-primary col-6 mx-auto">Добавить желание</button>
            </div>

            <ul>
                <div v-for="wish in wishes" :key="wish.id" class="col-8 col-lg-6 m-3 border border-1 border-dark px-3 py-5 rounded mx-auto">
                    {{ wish.title }} - {{ wish.description }}
                    <button v-if="isOwnUser()" @click="updateWish = { ...wish }" class="btn btn-warning">Редактировать</button>
                    <button v-if="isOwnUser()" @click="deleteWish(wish.id)" class="btn btn-danger">Удалить</button>
                </div>

                <div class="text-center" v-if="wishes.length === 0">Список желаний пуст</div>
            </ul>

            <div v-if="updateWish.id">
                <h3>Edit Wish</h3>
                <input v-model="updateWish.title" placeholder="Title" />
                <input v-model="updateWish.description" placeholder="Description" />
                <label>
                    Reserved:
                    <input v-model="updateWish.is_reserved" type="checkbox" />
                </label>
                <button @click="updateWishDetails" class="btn btn-primary">Update Wish</button>
            </div>
        </div>
        <button @click="logout" class="btn btn-danger logout-btn">Выйти</button>
    </div>
    `
};

// Define routes
const routes = [
    { path: '/auth', component: AuthPage },
    { path: '/:user_id', component: App }
];

// Create router
const router = createRouter({
    history: createWebHistory(),
    routes
});

// Mount the Vue app
createApp({ template: '<router-view />' }).use(router).mount('#app');
