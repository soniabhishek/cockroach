exports.up = function(knex, Promise) {
  	return Promise.all([
  		 knex
  			.schema
  			.table('work_flow', t => {
  			t.string('tag')
  					.notNullable()
  					.defaultTo('')
  			}),
         knex.raw(`UPDATE work_flow SET tag = feed_line.tag FROM feed_line, projects WHERE work_flow.project_id = projects.id AND feed_line.project_id = projects.id`),
 		 knex.raw('ALTER TABLE work_flow DROP CONSTRAINT work_flow_project_id_unique'),
 		 knex.raw('ALTER TABLE work_flow ADD CONSTRAINT work_flow_project_id_tag_unique UNIQUE (project_id,tag)')
  		])
  };

exports.down = function(knex, Promise) {
  	return Promise.all([
  	      knex
  		    .schema
  			.table('work_flow', t=> {
  				t.dropColumn('tag')
  			}),
  		  knex.raw('ALTER TABLE work_flow ADD CONSTRAINT work_flow_project_id_unique UNIQUE (project_id)')
  		])
};

