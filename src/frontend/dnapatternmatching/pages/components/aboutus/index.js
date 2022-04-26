import React from 'react';
import styles from './Aboutus.module.css';
import Subheading from '../subheading';
import ProfileCard from './components/profilecard';

const AboutUs = () => {
    return (
        <div className={styles.aboutContainer}>
            <Subheading 
                Text="About Us"
                Color="black"
            />
            <div className={styles.profileContainer}>
                <ProfileCard 
                    ImgPath="/aboutus/Rani.jpg"
                    ImgAlt="runnnayy"
                    Name="Maharani Ayu Putri Irawan"
                    NIM="13520019"
                />
                <ProfileCard 
                    ImgPath="/aboutus/Bryan.jpg"
                    ImgAlt="bryanbernigen"
                    Name="Bryan Bernigen"
                    NIM="13520034"
                />
                <ProfileCard 
                    ImgPath="/aboutus/Kyle.jpg"
                    ImgAlt="nk-kyle"
                    Name="Ng Kyle"
                    NIM="13520040"
                />
                {/* <div className={styles.profilesImage}>
                    <img src="aboutus/Rani.jpg" alt="runnnayy" className={styles.photo} />
                </div>
                <div className={styles.profilesImage}>
                    <img src="aboutus/Kyle.jpg" alt="bryanbernigen" className={styles.photo} />
                </div>
                <div className={styles.profilesImage}>
                    <img src="aboutus/Bryan.jpg" alt="nk-kyle" className={styles.photo} />
                </div> */}
            </div>
        </div>
    )
};

export default AboutUs;