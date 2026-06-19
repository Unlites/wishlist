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
                localStorage.setItem('username', this.user.name);
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
    <div class="min-vh-100 d-flex align-items-center justify-content-center" style="background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);">
        <div class="card shadow-lg border-0 rounded-4 p-3" style="max-width: 420px; width: 100%;">
            <div class="card-body">
                <h1 class="text-center fw-bold mb-1" style="color: #4f46e5;">Wishlist</h1>
                <p class="text-center text-muted mb-4">Список желаний</p>

                <div class="d-flex justify-content-center mb-4">
                    <div class="btn-group" role="group">
                        <button type="button" class="btn" :class="isRegister ? 'btn-outline-primary' : 'btn-primary'" @click="isRegister = false">Вход</button>
                        <button type="button" class="btn" :class="isRegister ? 'btn-primary' : 'btn-outline-primary'" @click="isRegister = true">Регистрация</button>
                    </div>
                </div>

                <div class="mb-3">
                    <label class="form-label text-muted small fw-semibold">Логин</label>
                    <input class="form-control form-control-lg" v-model="user.name" placeholder="Введите логин" autocomplete="username" />
                </div>
                <div class="mb-4">
                    <label class="form-label text-muted small fw-semibold">Пароль</label>
                    <input class="form-control form-control-lg" v-model="user.password" type="password" placeholder="Введите пароль" autocomplete="current-password" />
                </div>
                <button v-if="isRegister" class="btn btn-primary w-100 py-2 fw-semibold" @click="register">Зарегистрироваться</button>
                <button v-else class="btn btn-primary w-100 py-2 fw-semibold" @click="login">Войти</button>
            </div>
        </div>
    </div>
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
            localStorage.removeItem('username');
            this.$router.push('/auth');
        },
        isOwnUser() {
            return localStorage.getItem('user_id') == this.$route.params.user_id;
        },
        getLoggedInUsername() {
            return localStorage.getItem('username');
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
    <div class="app">
        <nav class="navbar navbar-expand navbar-dark shadow-sm" style="background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);">
            <div class="container">
                <span class="navbar-brand fw-bold">
                    <i class="bi bi-gift me-2"></i>Wishlist
                </span>
                <div class="d-flex align-items-center gap-3">
                    <span class="text-white-50 small d-none d-sm-inline">
                        <i class="bi bi-person-circle me-1"></i>{{ getLoggedInUsername() }}
                    </span>
                    <button @click="logout" class="btn btn-outline-light btn-sm">
                        <i class="bi bi-box-arrow-right me-1"></i>Выйти
                    </button>
                </div>
            </div>
        </nav>

        <div class="container py-4">
            <div class="card border-0 shadow-sm rounded-3 mb-4">
                <div class="card-body">
                    <div v-if="!isUserInfoUpdating" class="d-flex justify-content-between align-items-start">
                        <div>
                            <h5 class="card-title mb-1">
                                <i class="bi bi-person-circle me-2" style="color: #4f46e5;"></i>{{ user.username }}
                            </h5>
                            <p v-if="user.info" class="text-muted mb-0" style="white-space: pre-line">{{ user.info }}</p>
                            <p v-else class="text-muted fst-italic mb-0">Информация не добавлена</p>
                        </div>
                        <button v-if="isOwnUser()" class="btn btn-outline-secondary btn-sm flex-shrink-0" @click="isUserInfoUpdating = true">
                            <i class="bi bi-pencil"></i>
                            <span class="d-none d-sm-inline ms-1">{{ user.info ? 'Изменить' : 'Добавить' }}</span>
                        </button>
                    </div>
                    <form v-else @submit.prevent="updateUserInfo">
                        <h6 class="fw-semibold mb-2">{{ user.info ? 'Редактировать' : 'Добавить' }} информацию</h6>
                        <textarea @keydown.enter.exact.prevent="updateUserInfo" class="form-control" v-model="user.info" placeholder='Например, "Удобнее получить на OZON, мой пункт выдачи на Пушкина 36"' rows="3"></textarea>
                        <div class="mt-2 d-flex gap-2">
                            <button class="btn btn-primary btn-sm">
                                <i class="bi bi-check-lg me-1"></i>Сохранить
                            </button>
                            <button type="button" @click="isUserInfoUpdating = false" class="btn btn-outline-danger btn-sm">
                                <i class="bi bi-x-lg me-1"></i>Отмена
                            </button>
                        </div>
                    </form>
                </div>
            </div>

            <div v-if="isOwnUser()" class="card border-0 shadow-sm rounded-3 mb-4">
                <div class="card-body">
                    <h5 class="card-title fw-semibold mb-3">
                        <i class="bi bi-plus-circle me-2" style="color: #4f46e5;"></i>Новое желание
                    </h5>
                    <form @submit.prevent="addWish">
                        <div class="mb-2">
                            <input class="form-control" v-model="newWish.title" placeholder="Название" />
                        </div>
                        <div class="mb-2">
                            <textarea @keydown.enter.exact.prevent="$refs.addWishButton.click()" class="form-control" v-model="newWish.description" placeholder="Ссылка или описание" rows="2"></textarea>
                        </div>
                        <div class="mb-3">
                            <input type="number" min="1" step="any" class="form-control" v-model="newWish.price" placeholder="Цена (необязательно)" />
                        </div>
                        <button ref="addWishButton" class="btn btn-primary">
                            <i class="bi bi-plus-lg me-1"></i>Добавить желание
                        </button>
                    </form>
                </div>
            </div>

            <div v-if="wishes.length === 0" class="text-center py-5 empty-state">
                <i class="bi bi-inbox" style="font-size: 3rem; color: #adb5bd;"></i>
                <p class="text-muted fs-5 mt-2">Список желаний пуст</p>
            </div>

            <div class="row g-3" v-else>
                <div v-for="wish in wishes" :key="wish.id" class="col-12 col-md-6 col-lg-4">
                    <div class="card border-0 shadow-sm rounded-3 h-100 wish-card">
                        <div class="card-body d-flex flex-column">
                            <form v-if="wish.isUpdating" @submit.prevent="updateWishDetails">
                                <div class="mb-2">
                                    <input class="form-control form-control-sm" v-model="updateWish.title" placeholder="Название" />
                                </div>
                                <div class="mb-2">
                                    <textarea @keydown.enter.exact.prevent="updateWishDetails" class="form-control form-control-sm" v-model="updateWish.description" placeholder="Ссылка/описание" rows="2"></textarea>
                                </div>
                                <div class="mb-3">
                                    <input type="number" min="1" step="any" class="form-control form-control-sm" v-model="updateWish.price" placeholder="Цена" />
                                </div>
                                <div class="d-flex gap-2">
                                    <button class="btn btn-primary btn-sm">
                                        <i class="bi bi-check-lg me-1"></i>Сохранить
                                    </button>
                                    <button type="button" @click="stopWishUpdating(wish)" class="btn btn-outline-danger btn-sm">
                                        <i class="bi bi-x-lg me-1"></i>Отмена
                                    </button>
                                </div>
                            </form>
                            <div v-else class="d-flex flex-column h-100">
                                <div class="d-flex justify-content-between align-items-start mb-2">
                                    <h5 class="card-title fw-semibold text-break mb-0">{{ wish.title }}</h5>
                                </div>
                                <p class="card-text text-muted flex-grow-1 text-break mb-2" v-html="wish.description" style="white-space: pre-line"></p>
                                <div v-if="wish.price" class="mb-3">
                                    <span class="badge bg-success bg-opacity-10 text-success px-3 py-2 fs-6 fw-normal">
                                        <i class="bi bi-currency-ruble me-1"></i>{{ wish.price }}
                                    </span>
                                </div>
                                <div class="mt-auto">
                                    <div v-if="wish.is_reserved === true" class="d-flex flex-column gap-1">
                                        <button class="btn btn-success w-100" disabled>
                                            <i class="bi bi-check-circle-fill me-1"></i>Забронировано
                                        </button>
                                        <button v-if="isReservedByUser(wish)" @click="updateWishReserving(wish.id, false)" class="btn btn-outline-danger btn-sm w-100">
                                            <i class="bi bi-x-circle me-1"></i>Снять бронь
                                        </button>
                                    </div>
                                    <div v-else-if="wish.is_reserved === false">
                                        <button @click="updateWishReserving(wish.id, true)" class="btn btn-outline-primary w-100">
                                            <i class="bi bi-hand-index-thumb me-1"></i>Забронировать
                                        </button>
                                    </div>
                                    <div v-else class="d-flex gap-2">
                                        <button @click="startWishUpdating(wish)" class="btn btn-outline-secondary flex-grow-1">
                                            <i class="bi bi-pencil me-1"></i>Редактировать
                                        </button>
                                        <button @click="deleteWish(wish.id)" class="btn btn-outline-danger flex-grow-1">
                                            <i class="bi bi-trash me-1"></i>Удалить
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
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
