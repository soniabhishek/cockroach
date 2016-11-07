
exports.up = function(knex, Promise) {
  return Promise.all([
  knex.schema
    .createTable('feed_line_tags',t=>{
            t.increments()
            t.var('tag_name')
                .notNullable()
                .index()
            t.uuid('project_id')
                .notNullable()
                .index()
            t.uuid('work_flow_id')
            t.timestamp('created_at')
            t.timestamp('updated_at')
            t.unique('tag_name')
    }),
   knex.schema
    .table('work_flow', t=> {
    t.dropColumn('tag')
    })

  ])
};

exports.down = function(knex, Promise) {
 	return Promise.all([
  		 knex
  			.schema
  			.table('work_flow', t => {
  			t.string('tag')
  					.notNullable()
  					.defaultTo('')
  			})
  		 knex
  		    .schema
  		    .dropTableIfExists('feed_line_tags')
  		])
};
