import co from 'co'

exports.up = function(knex, Promise) {
    return co(function*() {
        yield knex.raw(`alter table crowdsourcing_flu_buffer add constraint crowdsourcing_flu_buffer_question_id_unique UNIQUE (question_id)`)
        yield knex.raw(`CREATE INDEX CONCURRENTLY crowdsourcing_flu_buffer_flu_id_index ON crowdsourcing_flu_buffer (flu_id)`)
    })
};

exports.down = function(knex, Promise) {
    return co(function*(){
        yield knex.raw(`ALTER TABLE crowdsourcing_flu_buffer DROP CONSTRAINT crowdsourcing_flu_buffer_question_id_unique`)
        yield knex.raw(`DROP INDEX CONCURRENTLY crowdsourcing_flu_buffer_flu_id_index`)
    })
};

exports.config = { transaction: false };