import React, { useEffect, useState } from 'react';
import { User, Notification } from '../models/models';

const NotificationsPage: React.FC<{ user: User }> = ({ user }) => {
    const [notifications, setNotifications] = useState<Notification[]>([]);

    useEffect(() => {
        // Fetch notifications for the user
        const fetchNotifications = async () => {
            try {
                const response = await fetch(`/api/notifications?userId=${user.id}`);
                const data = await response.json();
                setNotifications(data);
            } catch (error) {
                console.error('Error fetching notifications:', error);
            }
        };

        fetchNotifications();
    }, [user.id]);

    return (
        <div>
            <h1>Notifications for {user.username}</h1>
            <ul>
                {notifications.map(notification => (
                    <li key={notification.id}>{notification.message}</li>
                ))}
            </ul>
        </div>
    );
};

export default NotificationsPage;