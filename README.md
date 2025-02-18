# DigiEggs API

DigiEggs API is a simple RESTful API built with Go and the Gin framework, utilizing MongoDB for data storage. This API allows you to manage DigiEgg resources, where each DigiEgg represents an egg that is in a certain stage of incubation.

## Features
- **Create** a new DigiEgg
- **Read** all DigiEggs or a specific DigiEgg by ID
- **Update** a DigiEgg by ID
- **Delete** a DigiEgg by ID

## Technologies
- **Go**: The main programming language.
- **Gin**: A lightweight web framework for Go.
- **MongoDB**: A NoSQL database for storing the DigiEgg data.
- **dotenv**: To manage environment variables.

## Prerequisites

Before running this project, ensure that you have the following installed:

- [Go](https://golang.org/doc/install)
- [MongoDB](https://www.mongodb.com/try/download/community) (or use MongoDB Atlas for cloud hosting)
- [MongoDB Compass](https://www.mongodb.com/products/compass) (optional, for managing MongoDB visually)
- [Gin framework](https://github.com/gin-gonic/gin) (This is installed automatically by running `go mod`)

