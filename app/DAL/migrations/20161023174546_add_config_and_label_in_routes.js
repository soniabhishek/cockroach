
exports.up = function (knex, Promise) {
	return Promise.all([
		knex
			.schema
			.table('routes', t => {
				t.specificType('config', 'jsonb')
					.defaultTo('{}')
					.notNullable()
				t.string('label')
					.notNullable()
					.defaultTo('')
			})
	])
};

exports.down = function (knex, Promise) {
	return Promise.all([
		knex
			.schema
			.table('routes', t=> {
				t.dropColumn('config')
				t.dropColumn('label')
			})
	])
};
