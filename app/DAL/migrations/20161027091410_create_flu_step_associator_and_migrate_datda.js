import co from 'co'

/**
 * Dont un-comment it unless you know the consequences
 * @param knex
 * @param Promise
 */
exports.up = function(knex, Promise) {
	//return co(function*(){
	//
	//	yield knex.schema
	//		.createTable('sub_unit', t=> {
	//			t.increments()
	//			t.uuid('flu_id')
	//				.notNullable()
	//				.index()
	//			t.uuid('step_id')
	//				.notNullable()
	//				.index()
	//			t.integer('index')
	//				.notNullable()
	//			t.timestamp('created_at')
	//				.notNullable()
	//			t.timestamp('updated_at')
	//				.notNullable()
	//
	//			t.unique(['flu_id','step_id'])
	//			t.unique(['flu_id','index'])
	//		})
	//	yield knex.raw(`insert into sub_unit (flu_id, step_id, index, created_at, updated_at)
	//			SELECT id, step_id, 0,now(),now()
	//			FROM feed_line WHERE step_id is not NULL;`)
	//})
};

exports.down = function(knex, Promise) {
	//return co(function*(){
	//	yield knex.schema
	//		.dropTable('sub_unit')
	//})
};
