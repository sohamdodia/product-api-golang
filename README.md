# End Points

  * `/signup`
      * `POST /`
          * Creates a new user
          * **Requires**: `name, email, password`
          * **Accepts**: `name, email, password`
  * `/signin`
      * `POST /`
          * User signin
          * **Requires**: `email, password`
          * **Accepts**: `email, password`
  * `/product`
      * `GET /`
          * Get all products
          * **Requires**: No parameters
          * **Accepts**: No parameters
      * `GET /:id`
          * Get a single product
          * **Requires**: No parameters
          * **Accepts**: No parameters
      * `POST /`
          * Creates a new product
          * **Requires**: `name, oldPrice, newPrice`
          * **Accepts**: `name, oldPrice, newPrice`
      * `PUT /:id`
          * Updates a product.
          * **Requires**: No parameters
          * **Accepts**: `name, oldPrice, newPrice`
      * `DELETE /:id/`
          * Deletes a product
          * **Accepts**: No parameters
  