import co from 'co'

exports.up = function(knex, Promise) {
    return co(function*() {
        yield knex.raw(`CREATE INDEX CONCURRENTLY crowdsourcing_flu_buffer_flu_id_index ON crowdsourcing_flu_buffer (flu_id)`)
    })
};

exports.down = function(knex, Promise) {
    return co(function*(){
        yield knex.raw(`DROP INDEX CONCURRENTLY crowdsourcing_flu_buffer_flu_id_index`)
    })
};

exports.config = { transaction: false };
