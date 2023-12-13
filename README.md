# plspay
Split your group expenses easily with PlsPay

## Routes
| route             | methods | auth  | description                                                                                                   |
| :---------------- | :------ | :---: | :------------------------------------------------------------------------------------------------------------ |
| **/**             | GET     |       | Render a page to login, guest login or signup                                                                 |
| **/groups**       | GET     |   X   | Render a page with a table of groups for the logged-in user and a button to create new groups for that user   |
| **/groups**       | POST    |   X   | Endpoint to create a new group for the logged-in user                                                         |
| **/groups/:id**   | GET     |   x   | Render a page with #id group information and a table with expenses related to that group (add/delete expense) |
| **/users/signup** | GET     |       | Render a page with a form to signup a new user                                                                |
| **/users/signup** | POST    |       | Endpoint to signup a new user                                                                                 |
| **/users/login**  | GET     |       | Render a page with a form to login user                                                                       |
| **/users/login**  | POST    |       | Endpoint to login user                                                                                        |
| **/users/logout** | POST    |   x   | Endpoint to logout user                                                                                       |
| **/users/me**     | GET     |   x   | Render a page with logged-in user information                                                                 |