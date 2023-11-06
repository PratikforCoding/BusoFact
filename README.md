# BusoFact Project Readme

BusoFact is a GoLang-based backend server designed to help users find buses based on source, destination, or a specific bus number. This community-driven web server allows the general public to contribute by adding new buses or additional bus stoppages to existing routes. The system relies on MongoDB Atlas for data storage and incorporates JWT web tokens for authentication and security. The bcrypt package is used for password encryption.

## Technology Stack

- **Golang**: The primary programming language for building the server and handling backend operations.
- **MongoDB Atlas**: The database platform used to store information about buses, users, and other related data.
- **JWT (JSON Web Token)**: JWTs are employed for user authentication and authorization, ensuring secure access to the system.
- **bcrypt Package**: Used for hashing and encrypting user passwords to protect sensitive data.

## Features

### Role-Based Access

BusoFact implements a role-based access control system with two primary roles:

- **Consumer**: Consumers have limited capabilities, such as searching for buses and adding new buses or stoppages. They are not allowed to delete buses or users.

- **Admin**: Administrators have more extensive permissions, including the ability to delete buses and users. They have full control over the system. An "endpoint" called "make admin" allows a consumer to be promoted to an admin role.

### Authentication

Every user is required to log in to the system before they can add new buses or stoppages. User information, including their name, email address, and an encrypted password, is stored in the user collection of the MongoDB database.

- **Password Encryption**: User passwords are stored in an encrypted form to enhance security. This encryption is generated using the bcrypt package.

- **JWT Web Tokens**: JWTs are employed for user authentication. Users receive a JWT web token upon successful login, which is sent to the client-side via a cookie. This token is used to remember users, eliminating the need to log in every time they access the system. There are two types of tokens:

    - **Access Token**: Has a short lifespan for security purposes and is used for initial access to the system.
    
    - **Refresh Token**: Has a longer lifespan and is used to obtain new access tokens without requiring users to log in again.

### Containerization

BusoFact is designed for easy containerization, making it suitable for deployment in a microservices architecture.

Feel free to explore and contribute to the BusoFact project. If you have any questions or need assistance, please refer to the project documentation or contact the project maintainers.

For more information, visit the [BusoFact GitHub repository](https://github.com/PratikforCoding/BusoFact.git) and get involved in this community-driven project.

