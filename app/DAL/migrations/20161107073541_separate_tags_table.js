import co from 'co'

exports.up = function(knex, Promise) {
  return co(function*(){

    yield knex.schema
    .createTable('work_flow_tag_associators',t=>{
            t.increments()
            t.string('tag_name')
                .notNullable()
            t.uuid('work_flow_id')
                .notNullable()
                .index()
            t.uuid('project_id')
        .notNullable()
            t.timestamp('created_at')
            t.unique(['tag_name','project_id'])
            t.foreign('work_flow_id').references('work_flow.id')
            t.foreign('project_id').references('projects.id')
    })
    yield knex.raw(`
      INSERT INTO work_flow_tag_associators (tag_name, work_flow_id, project_id, created_at)
      SELECT tag,id,project_id,now() FROM work_flow;`)
    yield knex.schema
    .table('work_flow', t=> {
      t.dropColumn('tag')
    })

  })
};

exports.down = function(knex, Promise) {
 	return Promise.all([
  		 knex
  			.schema
  			.table('work_flow', t => {
  			t.string('tag')
  					.notNullable()
  					.defaultTo('')
  			}),
  		 knex
  		    .schema
  		    .dropTableIfExists('work_flow_tag_associators')
  		])
};
