exports.up = function(knex, Promise) {
  	return Promise.all([
  		 knex
  			.schema
  			.table('work_flow', t => {
  			t.string('tag')
  					.notNullable()
  					.defaultTo('')
  			}),
         knex.raw(`UPDATE work_flow SET tag = feed_line.tag FROM feed_line, projects WHERE work_flow.project_id = projects.id AND feed_line.project_id = projects.id`)

  		])
  };

exports.down = function(knex, Promise) {
  	return Promise.all([
  	      knex
  		    .schema
  			.table('work_flow', t=> {
  				t.dropColumn('tag')
  			})
  		])
};

