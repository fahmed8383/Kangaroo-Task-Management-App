# Kangaroo
 
The frontend of this project was generated with [Angular CLI](https://github.com/angular/angular-cli) version 9.1.6.

Kangaroo is a simple scheduling and task management app built using angular, golang for the back end, and postgresql for the database.

**Full support for mobile and touchscreen**

To build on new server:

1. Create a secrets.json file in /backend
2. Populate the file with 'recaptcha', 'dbuser', 'dbpassword', 'key', 'jwt', and 'gmailPass'
3. Run:
```
ng build --prod
```
4. Once the frontend build is complete, run:
```
docker-compose up -d --build
```
