import { createApp } from 'vue';
import axios from 'axios';
import { createRouter, createWebHistory } from 'vue-router';
import '../assets/styles.css';
import striptags from 'striptags';

const API_BASE_URL = '/api/v1';

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

const AuthPage = {
    data() {
        return {
            user: {
                name: '',
                password: '',
                token: ''
            },
            isRegister: false,
        };
    },
    methods: {
        async register() {
            try {
                await axios.post(`${API_BASE_URL}/users/register`, {
                    name: this.user.name,
                    password: this.user.password
                });
                this.login();
            } catch (error) {
                if (error.response && error.response.data) {
                    if (error.response.status === 400) {
                        if (error.response.data.includes('password: the length must be between 4 and 64')) {
                            alert('Пароль должен быть не менее 4 символов');
                        } else {
                            alert('Такой пользователь уже существует');
                        }
                    } else {
                        alert('Ошибка регистрации: ' + error.response.data);                        
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
                const redirectTo = this.$route.query.from || `/${tokenPayload.sub}`;
                this.$router.push(redirectTo);
            } catch (error) {
                if (error.response && error.response.data) {
                    if (error.response.status === 401) {
                        alert('Неверный логин или пароль'); 
                    } else if (error.response.status === 400) {
                        alert('Пароль должен быть не менее 4 символов');
                    } else {
                        alert('Ошибка аутентификации: ' + error.response.data);
                    }
                } else {
                    alert('Ошибка аутентификации, обратитесь к администратору');
                }
            }  
        }
    },
    template: `
    <form @submit.prevent class="auth p-3">
        <h1 class="text-center">Wishlist</h1>
        <hr>

        <div class="text-center m-5 row col-10 col-lg-4 mx-auto border border-2 border-dark px-3 py-5 rounded">
        <h3 class="text-center" v-if="!isRegister">
            Вход | <a href="#" class="link-primary" @click="isRegister = true">Регистрация</a>
        </h3>
        <h3 class="text-center" v-else>
            <a href="#" class="link-primary" @click="isRegister = fasle">Вход</a> | Регистрация
        </h3>

        <input class="mt-2" v-model="user.name" placeholder="Логин" />
            <input class="mt-2" v-model="user.password" type="password" placeholder="Пароль" />
            <div>
                <button v-if="isRegister" class="mt-3 btn btn-primary mx-auto px-5" @click="register">Зарегистрироваться</button>
                <button v-else class="mt-3 btn btn-primary mx-auto px-5" @click="login">Войти</button>
            </div>
        </div>
    </form>
    `
};

const App = {
    data() {
        return {
            user: {
                username: '',
                info: '',
            },
            isUserInfoUpdating: false,
            wishes: [],
            newWish: {
                title: '',
                description: '',
                price: null
            },
            updateWish: {
                id: null,
                title: '',
                description: '',
                price: null
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
                this.wishes = response.data.map(wish => {
                    wish.description = linkifyHtml(wish.description, { attributes: { target: '_blank' } });
                    return wish;
                });
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
                this.newWish = { title: '', description: '', price: null };
                await this.fetchWishes();
            } catch (error) {
                if (error.response && error.response.data) {
                    if (error.response.status === 400) {
                        if (error.response.data.includes('title: cannot be blank')) {
                            alert('Название не может быть пустым');
                        } else if (error.response.data.includes('title: the length must be between 1 and 300')) {
                            alert('Название не может быть больше 300 символов');
                        } else if (error.response.data.includes('description: the length must be between 1 and 5000')) {
                            alert('Описание не может быть больше 5000 символов');
                        }
                    } else {
                        alert('Ошибка добавления желания: ' + error.response.data);
                    }
                } else {
                    alert('Ошибка добавления желания');
                }
            }
        },
        async updateWishDetails() {
            try {
                if (this.updateWish.price === "") {
                    this.updateWish.price = null;
                }
                await api.put(`${API_BASE_URL}/users/${this.$route.params.user_id}/wishes/${this.updateWish.id}`, this.updateWish, {
                    headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
                });
                this.updateWish = { id: null, title: '', description: '', price: null };
                await this.fetchWishes();
            } catch (error) {
                if (error.response && error.response.data) {
                    if (error.response.status === 400) {
                        if (error.response.data.includes('title: cannot be blank')) {
                            alert('Название не может быть пустым');
                        } else if (error.response.data.includes('title: the length must be between 1 and 300')) {
                            alert('Название не может быть больше 300 символов');
                        } else if (error.response.data.includes('description: the length must be between 1 and 5000')) {
                            alert('Описание не может быть больше 5000 символов');
                        }
                    } else {
                        alert('Ошибка обновления желания: ' + error.response.data);                        
                    }
                } else {
                    alert('Ошибка обновления желания'); 
                }
            }
        },
        async deleteWish(wishId) {
            if (!confirm('Вы действительно хотите удалить желание?')) {
                return;
            }
            try {
                await api.delete(`${API_BASE_URL}/users/${this.$route.params.user_id}/wishes/${wishId}`, {
                    headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
                });
                await this.fetchWishes();
            } catch (error) {
                alert('Ошибка удаления желания');
            }
        },
        async updateWishReserving(wishId, isReserved) {
            try {
                await api.put(`${API_BASE_URL}/users/${this.$route.params.user_id}/wishes/${wishId}/update-reserving`, 
                    {
                        is_reserved: isReserved
                    }, {
                    headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
                });
                await this.fetchWishes();
            } catch (error) {
                alert('Ошибка резервирования желания');
            }
        },
        async fetchUser() {
            try {
                const response = await api.get(`${API_BASE_URL}/users/${this.$route.params.user_id}`, {
                    headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
                });
                this.user.username = response.data.name;
                this.user.info = response.data.info;
            } catch (error) {
                alert('Ошибка загрузки пользователя');
            }
        },
        async updateUserInfo() {
            try {
                await api.put(`${API_BASE_URL}/users/${this.$route.params.user_id}/info`,
                    { info: this.user.info },
                    { headers: { Authorization: `Bearer ${localStorage.getItem('token')}` } },
                );
                await this.fetchUser();
                this.isUserInfoUpdating = false;
            } catch (error) {
                if (error.response && error.response.data) {
                    if (error.response.status === 400) {
                        if (error.response.data.includes('info: the length must be between 0 and 3000')) {
                            alert('Информация не может быть больше 3000 символов');
                        }
                    } else {
                        alert('Ошибка изменения информации для пользователей: ' + error.response.data);                        
                    }
                } else {
                    alert('Ошибка изменения информации для пользователей'); 
                }
            }
        },
        logout() {
            localStorage.removeItem('token');
            localStorage.removeItem('user_id');
            this.$router.push('/auth');
        },
        isOwnUser() {
            return localStorage.getItem('user_id') == this.$route.params.user_id;
        },
        isReservedByUser(wish) {
            if (!wish.reserved_by) {
                return false;
            }

            return wish.reserved_by == localStorage.getItem('user_id');
        },
        startWishUpdating(wish) {
            this.updateWish = { ...wish }; 
            this.updateWish.description = striptags(wish.description);            
            wish.isUpdating = true;
        },
        stopWishUpdating(wish) {
            this.updateWish = { id: null, title: '', description: '', price: null };
            wish.isUpdating = false;
        }
    },
    mounted() {
        this.fetchUser();
        this.fetchWishes();
    },
    template: `
    <div class="app p-3">
        <div class="text-center">
            <h1>Wishlist {{ user.username }}</h1>

            <div class="mt-3" v-if="!isUserInfoUpdating" style="white-space: pre-line">{{ user.info }}</div>
            
            <button v-if="isOwnUser() && !isUserInfoUpdating" class="btn btn-outline-primary mt-1" @click="isUserInfoUpdating = true"><span v-if="!user.info">Добавить</span><span v-else>Изменить</span> информацию для пользователей</button>
            
            <form @submit.prevent="updateUserInfo" v-if="isUserInfoUpdating" class="text-center row col-10 col-md-6 col-lg-6 mx-auto">
                <textarea @keydown.enter.exact.prevent="updateUserInfo" v-model="user.info" placeholder='Например, "Удобнее получить на OZON, мой пункт выдачи на Пушкина 36"' class="mt-3" />
                <div class="mt-2">
                    <button class="m-1 btn btn-primary">Сохранить</button>
                    <button type="button" @click="isUserInfoUpdating = false" class="btn btn-outline-danger">Отмена</button>
                </div>
            </form>
            
        </div>
        <hr>

        <div>
            <form @submit.prevent="addWish" v-if="isOwnUser()" class="text-center m-5 row col-10 col-md-8 col-lg-6 col-xl-4 mx-auto border border-2 border-dark px-3 py-5 rounded">
                <h3>Новое желание</h3>
                <input class="m-1" v-model="newWish.title" placeholder="Название" />
                <textarea @keydown.enter.exact.prevent="$refs.addWishButton.click()" class="m-1" v-model="newWish.description" placeholder="Ссылка/описание" />
                <input type="number" min="1" step="any" class="m-1" v-model="newWish.price" placeholder="Цена (необязательно)" />
                <button ref="addWishButton" class="mt-3 btn btn-primary col-lg-6 mx-auto">Добавить желание</button>
            </form>

            <div>
                <div v-for="wish in wishes" :key="wish.id" class="text-center col-10 col-md-10 col-lg-8 col-xl-6 m-3 border border-1 border-dark px-3 py-5 rounded mx-auto">
                    <form @submit.prevent="updateWishDetails" v-if="wish.isUpdating" class="row text-center col-10 col-lg-8 mx-auto">
                        <input class="mt-2" v-model="updateWish.title" placeholder="Название" />
                        <textarea @keydown.enter.exact.prevent="updateWishDetails" class="mt-2" v-model="updateWish.description" placeholder="Ссылка/описание" />
                        <input type="number" min="1" step="any" class="mt-2" v-model="updateWish.price" placeholder="Цена (необязательно)" />
                        <div class="mt-2">
                            <button class="m-1 btn btn-primary">Сохранить</button>
                            <button type="button" @click="stopWishUpdating(wish)" class="btn btn-outline-danger">Отмена</button>
                        </div>
                    </form>
                    <div v-else>
                        <h3 class="text-break">{{ wish.title }}</h3>
                        <p class="text-break" v-html="wish.description"></p>
                        <b v-if="wish.price">Цена: {{ wish.price }} руб.</b>
                        <div class="mt-3">
                            <div v-if="wish.is_reserved === true" class="d-flex justify-content-center align-items-center flex-column">
                                <button class="btn btn-primary px-5" disabled>Забронировано</button>
                                <button v-if=isReservedByUser(wish) @click="updateWishReserving(wish.id, false)" class="btn btn-outline-danger px-2 py-1 mt-1">Снять бронь</button>
                            </div>
                            <div v-else-if="wish.is_reserved === false">
                                <button @click="updateWishReserving(wish.id, true)" class="btn btn-outline-primary">Забронировать</button>
                            </div>
                            <div v-else>
                                <button @click="startWishUpdating(wish)" class="btn btn-outline-dark m-1">Редактировать</button>
                                <button @click="deleteWish(wish.id)" class="btn btn-outline-danger">Удалить</button>
                            </div>
                        </div>
                    </div>
                </div>

                <h2 class="text-center mt-5" v-if="wishes.length === 0">Список желаний пуст</h2>
            </div>
        </div>
        <button @click="logout" class="btn btn-danger logout-btn">Выйти</button>
    </div>
    `
};

// Define routes
const routes = [
    { path: '/auth', component: AuthPage },
    { path: '/:user_id', component: App, meta: { requiresAuth: true } },
    { 
        path: '/', 
        redirect: () => {
            const userId = localStorage.getItem('user_id');
            return userId ? `/${userId}` : '/auth';
        }
    },
];

// Create router
const router = createRouter({
    history: createWebHistory(),
    routes
});

router.beforeEach((to, from, next) => {    
    if (to.meta.requiresAuth && !localStorage.getItem('token')) {
        next({ path: '/auth', query: { from: to.params.user_id } });
    } else {
        next();
    }
});

const app = createApp({ template: '<router-view />' });
app.use(router);
app.mount('#app');
