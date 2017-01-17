
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
            .dropTableIfExists('logic_gate_formula')

    ])
};

exports.down = function(knex, Promise) {
  
};
