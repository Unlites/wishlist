# Wishlist

Service for managing and sharing your wishes - [super-wishlist.ru](https://super-wishlist.ru)

## Install

```shell
git clone https://github.com/Unlites/wishlist
```

After cloning the project, create .env file based on .env.sample with overriding parameters if necessary.

Make sure that Docker service is currently running on your system.

### Local

Service can be run locally for development purposes:

```shell
docker-compose -f docker-compose.local.yml up -d
cd ui
npm run dev
```

Now you can access the service at http://localhost:5173/

### Production ready

To use the service as a production ready application, you have to configure ```certbot``` first with official script following the volumes defined in ```docker-compose.yml```. Then:

```shell
docker-compose up -d
```

## Features

Any authenticated user can build his own wishlist and share it with others. 

Users can also see other users wishlists and mark chosen wishes as reserved for a gift.

Owner of the wishlist cannot see reserving status of his wishes to keep a surprise part.