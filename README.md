# Test task for Effective mobile

This service is used to store person's info which is enriched with `age`, `gender` and `nationality` info by sending requsts to third-party api's.  
Implement all requsted rest methods: create, get with filters, update and delete person.  
Pagination implemented using (`person_id`, `created_at`)-way. Get persons request accept filters from Query params. All the params are similar as in `Person` struct and also `created_at`, `limit`. `limit` param is obligatory, other ones are not.  
Updating person implemented using pointers in request struct to check it for nil.  
The main technologies are:  
- `chi-router` for routing;
- `pgx` for postgres driver;
- `sqlx` as add-on `database/sql` package to work with database;
- `slog` as logger;
- `goose` as migration tool for database;
- `godotenv` to work with environment variables;

`.env` file was added to gitignore. Configurable variables are:
- `SERVER_PORT`;
- `SERVER_READ_TIMEOUT`;
- `SERVER_WRITE_TIMEOUT`;
- `DATABASE_URL`;
- `AGE_BASE_URL` is url for third-party api to find person age;
- `GENDER_BASE_URL` is url for third-party api to find person gender;
- `NATIONALITY_BASE_URL` is url for third-party api to find person nationality;

Implement graceful shutdown. Add debug, info and error logger.
