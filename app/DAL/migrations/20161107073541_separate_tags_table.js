
exports.up = function(knex, Promise) {
  return Promise.all([
  knex.schema
    .createTable('work_flow_tag_associators',t=>{
            t.increments()
            t.string('tag_name')
                .notNullable()
                .index()
            t.uuid('work_flow_id')
            t.timestamp('created_at')
            t.unique('tag_name')
            t.foreign('work_flow_id').references('work_flow.id')            
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
  			}),
  		 knex
  		    .schema
  		    .dropTableIfExists('work_flow_tag_associators')
  		])
};
