
exports.up = function(knex, Promise) {
  return knex.raw('INSERT INTO step_type VALUES (11, ?);', 'EMPTY_STEP')
};

exports.down = function(knex, Promise) {
  return knex.raw('DELETE FROM step_type where name = ?;', 'EMPTY_STEP')
};
