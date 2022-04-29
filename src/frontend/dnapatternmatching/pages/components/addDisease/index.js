import React, { useState } from 'react';
import styles from './AddDisease.module.css';
import Subheading from '../subheading';

const AddDisease = () => {
    const [namapenyakit, setNamapenyakit] = useState('');
    const [file, setFile] = useState(null);
    const [rantaidna, setRantaiDNA] = useState('');

    const handleChangeNamapenyakit = (e) => {
        setNamapenyakit(e.target.value);
        console.log(namapenyakit);
    }

    const handleChangeFile = (event) => { 
        event.preventDefault();
        setFile(file);
        
        const reader = new FileReader();
        reader.onload = (event) => {
            const text = event.target.result;
            console.log(text);
            setRantaiDNA(text);
        }
        if (event.target.files[0]) {
            reader.readAsText(event.target.files[0]);
        }
        console.log(rantaidna)
    };

    const onSubmitForm = async (e) => {
        e.preventDefault();
        console.log("tombol dipencet");
        fetch('http://localhost:8080/api/submitdisease' + '/' + namapenyakit + '/' + rantaidna, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                mode: 'no-cors'
        })
        .then(json => { 
            console.log(json);
            setNamapenyakit('');
            setRantaiDNA('');
            setFile(null);
        });
    }

    return (
        <div className={styles.addDiseaseContainer}>
            <Subheading 
                Text="Add Disease Data"
                Color="black"
            />
            <div className={styles.formDiseaseContainer}>
                <form method="post" onSubmit={onSubmitForm} className={styles.formCt}>
                    <label className={styles.label} >Disease Name: </label>
                    <input type="text" required value={namapenyakit} onChange={handleChangeNamapenyakit} className={styles.inputText} />
                    <label className={styles.label} >DNA Sequence: </label>
                    <input type="file" value={file} onChange={handleChangeFile} className={styles.inputFile} />
                    <button type="submit" className={styles.submitButton} >Submit</button>
                </form>
            </div>
        </div>
    )
};

export default AddDisease;