import co from 'co'

exports.up = function(knex, Promise) {
    return co(function * () {
        yield knex.schema.table('routes', t=> {t.dropColumn('logic_gate_id')})
        yield knex.schema.dropTableIfExists('logic_gate')
        yield knex.schema.dropTableIfExists('logic_gate_formula')
        yield knex.raw(`update routes set config = concat('{"input_template":[',r2.config, ']}')::jsonb from routes as r1 join routes as r2 on r1.id = r2.id`)
    })
};



exports.down = function(knex, Promise) {
    return co(function * () {

        yield knex.schema
                .createTable('logic_gate_formula',t => {
                    t.integer('id')
                    .primary()
                t.string('name')
        })

        yield knex.schema
        	    .createTable('logic_gate',t=>{
                    t.uuid('id')
                    .primary()
                t.specificType('input_template','jsonb')
                t.integer('formula')
                    .notNullable()
                    .references('id')
                    .inTable('logic_gate_formula')
        	    })

        yield knex.schema
                .table('routes', t => {
                    t.uuid('logic_gate_id')
                    .references('id')
                    .inTable('logic_gate')
                    .onDelete('CASCADE')
                })
    })
};
