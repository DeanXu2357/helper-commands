package tpl

func MigrationCreateTemplate() []byte {
	return []byte(`
const _ = require('lodash');
const { logger } = require('.');
const currentVersion = '{{ .Version }}';

module.exports.Execute = async (db, aql, logger, env) => {
    try {
        const defaultCol = db.collection('{{ .CollectionName }}');

        const hasDefaultCol = await defaultCol.exists();


        // ensure collections
        if (!hasDefaultCol) {
            await db.createCollection('{{ .CollectionName }}')
        }

        const cursor = await db.query(aql` + "`" + `
	FOR v IN Versions
	FILTER !IS_NULL(v.migrationNumber) AND
	v.migrationNumber == ${currentVersion}
	RETURN v._key
	` + "`" + `);
        const currentMigration = await cursor.next();
        if (!_.isUndefined(currentMigration)) {
            return;
        }

        // update the schema version to {{ .Version }}
        logger.info('update the schema version the {{ .Version }}');
        await db.collection('Versions').save(
            {
                migrationNumber: currentVersion,
                createdAt: Date.now(),
            },
            { returnNew: true },
        );
        logger.info('*** migration successfully execution ***');
        return;
    } catch (error) {
        console.log('shit gg :\n', error);
    }
}

module.exports.Rollback = async (db, aql, logger, env) => {
    try {
        const versionsCol = db.collection('Versions');
        const cursor = await db.query(aql` + "`" + `
	FOR v IN Versions
	FILTER !IS_NULL(v.migrationNumber) AND
	v.migrationNumber == ${currentVersion}
	RETURN v._key
	` + "`" + `);
        const currentMigration = await cursor.next();

        await dropCollection(db, '{{ .CollectionName }}');

        // 移除 Migration
        await versionsCol.remove(currentMigration);

        logger.info('*** rollback successfully execution ***');
        return;
    } catch (error) {
        console.log('Rollback Error:\n', error);
    }
}

const dropCollection = async (db, colName) => {
    const collection = db.collection(colName);
    const hasCol = await collection.exists();

    if (!hasCol) {
        logger.info('Collection %s not exists', colName);
    } else {
        collection.drop();
        logger.info('Drop %s collection success', colName);
    }
    return;
};
`)
}
