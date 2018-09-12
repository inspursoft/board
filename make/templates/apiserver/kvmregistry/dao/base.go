package dao

import (
	"database/sql"
	"log"
	"strconv"
)

func GetDBConn() (db *sql.DB) {
	db, err := sql.Open("sqlite3", "file:registry.db?journal_mode=WAL")
	if err != nil {
		log.Fatalf("Failed to get connection to DB: %+v", err)
	}
	return
}

func InitDB() (err error) {
	_, err = GetDBConn().Exec(`
		create table if not exists registry
		 (kvm_name text, job_name text, build_id, primary key(kvm_name, job_name));
		create table if not exists kvm
		 (kvm_name text, affinity text, primary key(kvm_name, affinity));
	`)
	return
}

func AddOrUpdateRegistry(kvmName, jobName string) (err error) {
	ptmt, err := GetDBConn().Prepare(`
		insert or replace into registry (kvm_name, job_name)
		 values (?, ?);
	`)
	_, err = ptmt.Exec(kvmName, jobName)
	return
}

func UpdateRegistry(buildID, kvmName, jobName string) (err error) {
	ptmt, err := GetDBConn().Prepare(`
	  update registry set build_id = ?
	    where kvm_name = ? and job_name = ?;
	`)
	_, err = ptmt.Exec(buildID, kvmName, jobName)
	return
}

func DeleteLastAffinity() (err error) {
	_, err = GetDBConn().Exec(`delete from kvm;`)
	return
}

func DeleteRelationalJob(jobName string, buildID string) (err error) {
	ptmt, err := GetDBConn().Prepare(`
		delete from registry where job_name = ? and build_id = ?; 
	`)
	_, err = ptmt.Exec(jobName, buildID)
	return
}

func DeleteRegistryWithJob(kvmName string, jobName string) (err error) {
	ptmt, err := GetDBConn().Prepare(`
		delete from registry where kvm_name = ? and job_name = ?; 
	`)
	_, err = ptmt.Exec(kvmName, jobName)
	return
}

func AddOrUpdateKVM(kvmName string, affinity string) (err error) {
	ptmt, err := GetDBConn().Prepare(`
		insert or replace into kvm (kvm_name, affinity)
		 values (?, ?);
	`)
	_, err = ptmt.Exec(kvmName, affinity)
	return
}

func GetKVMByJob(jobName string) (kvm string, err error) {
	rs, err := GetDBConn().Query(`
		select r.kvm_name kvm_name
			from registry r left join kvm k
				on k.kvm_name = r.kvm_name
			where r.job_name = ?;
	`, jobName)
	if err != nil {
		log.Printf("Failed to get KVM by job: %s, error: %+v\n", jobName, err)
		return "", err
	}
	for rs.Next() {
		err = rs.Scan(&kvm)
		return
	}
	return
}

func GetKVMStatus() (kvmList []map[string]string, err error) {
	rs, err := GetDBConn().Query(`
		select k.kvm_name kvm_name, k.affinity affinity, 
		  r.job_name job_name, r.build_id build_id,
		  case when r.kvm_name is null then 0 else 1 end is_allocated
			from kvm k left join registry r
			on k.kvm_name = r.kvm_name;
	`)
	if err != nil {
		log.Printf("Failed to get KVM status, error: %+v\n", err)
		return nil, err
	}
	kvmList = []map[string]string{}
	for rs.Next() {
		var kvm string
		var affinity string
		var jobName sql.NullString
		var buildID sql.NullString
		var isAllocated int
		err = rs.Scan(&kvm, &affinity, &jobName, &buildID, &isAllocated)
		item := map[string]string{
			kvm:            affinity,
			"job_name":     jobName.String,
			"build_id":     buildID.String,
			"is_allocated": strconv.Itoa(isAllocated),
		}
		kvmList = append(kvmList, item)
	}
	return
}

func GetAvailableKVMWithAffinity(affinity string) (registryList []map[string]string, err error) {
	rs, err := GetDBConn().Query(`
		select k.kvm_name kvm_name, r.job_name job_name
			from kvm k left join registry r
			on k.kvm_name = r.kvm_name
			where r.kvm_name is null 
			   and r.build_id is null
			   and k.affinity = ?;
	`, affinity)
	if err != nil {
		log.Printf("Failed to get available KVM with affinity, error: %+v\n", err)
		return nil, err
	}
	registryList = []map[string]string{}
	for rs.Next() {
		var kvmName string
		var jobName sql.NullString
		err = rs.Scan(&kvmName, &jobName)
		registry := make(map[string]string)
		registry["kvm_name"] = kvmName
		registry["job_name"] = jobName.String
		registryList = append(registryList, registry)
	}
	return
}
