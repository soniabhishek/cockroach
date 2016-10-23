var Config = require('config-js');
var config = new Config('./app/config/##.json');

var connection = config.get('postgres')
connection.user = connection.username

var knexConfig = {
	client: 'postgresql',
	connection: connection,
	migrations: {
		directory: `${__dirname}/app/DAL/migrations`,
		tableName: 'angel_knex_migrations'
	},
	seeds: {
		directory: './app/db/seeds'
	},
	pool: {
		min: 5,
		max: 50
	},
	acquireConnectionTimeout: 60000
}

module.exports = {
	development: knexConfig,
	staging: knexConfig,
	production: knexConfig
}