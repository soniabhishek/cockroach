
exports.up = function(knex, Promise) {
    return Promise.all([
        knex
            .schema
            .table('routes', t=> {
            t.dropColumn('logic_gate_id')}),
        knex
            .schema
            .dropTableIfExists('logic_gate'),
        knex
            .schema
            .dropTableIfExists('logic_gate_formula'),
        knex
            .raw(`update routes set config = concat('{"input_template":[',r2.config, ']}')::jsonb from routes as r1 join routes as r2 on r1.id = r2.id`)

    ])
};

exports.down = function(knex, Promise) {
};
