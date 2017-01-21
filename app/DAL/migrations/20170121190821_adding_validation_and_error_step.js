import co from 'co'

exports.up = function(knex, Promise) {
    return co(function*() {
        yield knex.raw(`insert into step_type values (12,'VALIDATION_STEP')`)
    })

};

exports.down = function(knex, Promise) {
    return co(function*() {
        yield knex.raw(`delete from step_type where id = 12`)
    })
};
