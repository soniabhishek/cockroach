import co from 'co'

exports.up = function(knex, Promise) {
	return co(function*() {
		yield knex
			.schema
			.table('clients', t => {
				t.string('name')
					.notNullable()
					.defaultTo('')
			})
		yield knex.raw(`update clients as c set name=u.username from users as u where c.user_id=u.id`)

		yield knex.raw(`alter table clients add constraint name_unique unique(name)`)
	})
};

exports.down = function(knex, Promise) {
	return co(function*() {
		yield knex.raw(`alter table clients drop column name`)
	})
};
