-- 1.1 -- число событий по дням
-- группируем по дням, для кааждого дня считаем число событий
SELECT COUNT() FROM ads_data GROUP BY date;


-- 1.2 -- число показов по дням
-- считаем, что event = 1, это показ
-- a event = 2, это клик
SELECT COUNT() FROM ads_data
  GROUP BY date
  HAVING event = 1;


-- 1.3 -- число кликов по дням
SELECT COUNT() FROM ads_data
  GROUP BY date
  HAVING event = 2;

-- 1.4 -- число уникальных объявлений по дням
-- для подсчет уникальных объявлений, используем DISTINCT
-- для идентификатора объявлений
SELECT COUNT(DISTINCT(ad_id)) FROM ads_data GROUP BY date;

-- 1.5 -- число уникальных кампаний по дням
-- для подсчет уникальных кампаний, используем DISTINCT
-- для идентификатора кампании
SELECT COUNT(DISTINCT(campaign_union_id)) FROM ads_data GROUP BY date;

-- 2 -- объявления, по которым показ произошел после клика
-- первым запросом clicks получаем список объявленией клики, со временем клика
-- вторым запросом объединяем показы с соответствующими им кликами
-- при этом оставляем только те объявления, у которых показ по времени идет
-- после клика
WITH clicks as (
  SELECT ad_id, time
  FROM ads_data
  WHERE event_id = 1
)
SELECT DISTINCT(ad_id)
  FROM ads_data JOIN clicks ON ads_data.ad_id = clicks.ad_id
  WHERE event_id = 2 and ads_data.time > clicks.time;
