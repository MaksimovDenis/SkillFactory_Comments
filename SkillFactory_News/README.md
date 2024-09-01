# news_agregator

# Запуск  

1. Сколнировать репозиторий:
```bash   
git clone https://github.com/MaksimovDenis/news_agregator.git
```

2. Перейти в директорию проекта (если Вы не в ней):  
```bash    
cd news_agregator 
```

3. Поднимет базу данных:  
```bash      
make up 
```

4. Локальная установка утилиты миграции:  
```bash      
make install-deps 
```

5. Накатить миграции:  
```bash      
make migration-up 
```

6. Запустить сервис:  
```bash      
make run 
```

 # Результат запроса /feeds/10

 - Создать пост    
![image](https://github.com/user-attachments/assets/29c61d7a-6378-4784-9404-9c080f474230)
