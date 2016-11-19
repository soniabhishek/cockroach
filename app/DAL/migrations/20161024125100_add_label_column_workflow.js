exports.up = function(knex, Promise) {
  	return Promise.all([
  		 knex
  			.schema
  			.table('work_flow', t => {
  			t.string('label')
  					.notNullable()
  					.defaultTo('')
  			}),
  		])
  };

exports.down = function(knex, Promise) {
  	return Promise.all([
  	      knex
  		    .schema
  			.table('work_flow', t=> {
  				t.dropColumn('label')
  			}),
  		])
};

