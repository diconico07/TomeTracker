# TomeTracker

TomeTracker is a full-stack application designed to help users manage their book series collections. It is meant to be self-hosted.

## Table of Contents

- [Features](#features)
- [Technologies Used](#technologies-used)
- [App Architecture](#app-architecture)

## Features

- **Add Books Series**: Easily add new book series to your tracked series
- **View Books**: Browse books from these series and mark whether you own them or not.
- **View Release Planning**: Keep an eye on announced books and see when they are to be released.

## Technologies Used

### Backend

- **MariaDB**: Relational database for storing book data.
- **Gin Gonic**: For the API
- **Ent**: For Database
- **goquery**: To Fetch data from editors' websites

### Frontend

- **Vue.js**: JavaScript library for building user interfaces.
- **Vuetify**: Material Design Component Library
- **Axios**: Promise-based HTTP client for the browser and Node.js.

## App Architecture

The application is based around multiple components:

- **TomeStore**: the backend long-running application that interacts with the database and serves an API for all other
- **UI**: The Vue frontend that calls the API
- **TomeFetchers**: These are meant to run as periodic jobs, they are responsible for fetching the actual books data from the editor's website, currently only a YenPress one is available.
