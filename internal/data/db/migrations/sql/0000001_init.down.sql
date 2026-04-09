DROP TRIGGER IF EXISTS trg_user_files_after_delete ON file_nodes;
DROP FUNCTION IF EXISTS trg_user_files_hard_delete();

DROP INDEX IF EXISTS uniq_active_sub;
DROP TABLE IF EXISTS user_subscriptions;
DROP TABLE IF EXISTS plans;

DROP TABLE IF EXISTS share_links;
DROP TABLE IF EXISTS file_nodes;
DROP TABLE IF EXISTS physical_files;

DROP FUNCTION IF EXISTS uuidv7();
