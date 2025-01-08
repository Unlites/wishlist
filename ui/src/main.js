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
                this.error = 'Registration failed';
            }
        },
        async login() {
            try {
                const response = await axios.post(`${API_BASE_URL}/users/login`, {
                    name: this.user.name,
                    password: this.user.password
                });
                this.user.token = response.data.token;
                localStorage.setItem('token', this.user.token);
                this.$router.push('/');
            } catch (error) {
                this.error = 'Login failed';
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
                this.error = 'Failed to fetch wishes';
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
                this.error = 'Failed to add wish';
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
                this.error = 'Failed to update wish';
            }
        },
        async deleteWish(wishId) {
            try {
                await api.delete(`${API_BASE_URL}/users/${this.$route.params.user_id}/wishes/${wishId}`, {
                    headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
                });
                await this.fetchWishes();
            } catch (error) {
                this.error = 'Failed to delete wish';
            }
        },
        logout() {
            localStorage.removeItem('token');
            this.$router.push('/auth');
        }
    },
    mounted() {
        this.fetchWishes();
    },
    template: `
    <div class="app p-3">
        <h1 class="text-center">Wishlist</h1>

        <div v-if="error" style="color: red;">{{ error }}</div>

        <div>
            <h2>Wishes</h2>

            <div class="text-center m-3 row">
                <h3>Add Wish</h3>
                <input v-model="newWish.title" placeholder="Название" />
                <input v-model="newWish.description" placeholder="Ссылка/описание" />
                <button @click="addWish" class="btn btn-primary col-3 mx-auto">Добавить желание</button>
            </div>

            <ul>
                <li v-for="wish in wishes" :key="wish.id">
                    {{ wish.title }} - {{ wish.description }}
                    <button @click="updateWish = { ...wish }" class="btn btn-warning">Редактировать</button>
                    <button @click="deleteWish(wish.id)" class="btn btn-danger">Удалить</button>
                </li>
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
        <button @click="logout" class="btn btn-danger">Logout</button>
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
