import React, { useState, useRef } from 'react';
import styles from './AddDisease.module.css';
import Subheading from '../subheading';

const AddDisease = () => {
    const [namapenyakit, setNamapenyakit] = useState('');
    const [file, setFile] = useState(null);
    const [rantaidna, setRantaiDNA] = useState('');

    const onSubmitForm = async () => {
        const response = await fetch('localhost:8080', {
            method: 'POST',
            body: JSON.stringify({
                namapenyakit,
                rantaidna,
            }),
            headers: {
                'Content-Type': 'application/json',
            },
        });
        const data = await response.json();
        console.log(data);
    }

    const handleChangeName = (event) => {
        setNamapenyakit(event.target.value);
        console.log(namapenyakit);
    };

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

    return (
        <div className={styles.addDiseaseContainer}>
            <Subheading 
                Text="Add Disease Data"
                Color="black"
            />
            <div className={styles.formDiseaseContainer}>
                <form action="/api/submitdisease/" method="post" onSubmit={onSubmitForm} className={styles.formCt}>
                    <label className={styles.label} >Disease Name: </label>
                    <input type="text" required value={namapenyakit} onChange={handleChangeName} className={styles.inputText} />
                    <label className={styles.label} >DNA Sequence: </label>
                    <input type="file" value={file} onChange={handleChangeFile} className={styles.inputFile} />
                    <button type="submit" className={styles.submitButton} >Submit</button>
                </form>
            </div>
        </div>
    )
};

export default AddDisease;

// import axios from 'axios';
// import { useForm } from 'react-hook-form';

    // const { register, handleSubmit, errors, reset } = useForm();
    // async function onSubmitForm(values) {
    //     let config = {
    //         method: 'post',
    //         url: `${process.env.NEXT_PUBLIC_API_URL}/api/contact`,
    //         headers: {
    //             'Content-Type': 'application/json',
    //         },
    //         data: values,
    //     };

    //     try {
    //         const response = await axios(config);
    //         console.log(response);
    //         if (response.status == 200) {
    //             reset();
    //         }
    //     } catch (err) {}
    // }
{/* <div className={styles.formDiseaseContainer}>
                <form
                    onSubmit={handleSubmit(onSubmitForm)} 
                    action="/api/new" method="post" className={styles.formCt}>
                    <label for="roll" className={styles.label} >Disease Name: </label>
                    <input type="text" ref={register} required className={styles.inputText} />
                    <label for="name" className={styles.label} >DNA Sequence: </label>
                    <input name="logo" ref={register} type="file" className={styles.inputFile} />
                    <button type="submit" className={styles.submitButton} >Submit</button>
                </form>
            </div> */}